package device

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/txchat/dtalk/api/proto/auth"
	"github.com/txchat/dtalk/api/proto/chat"
	"github.com/txchat/dtalk/api/proto/message"
	protoutil "github.com/txchat/im-util/internal/proto"
	"github.com/txchat/im-util/pkg/net"
	xhttp "github.com/txchat/im-util/pkg/net/http"
	"github.com/txchat/im-util/pkg/net/tcp"
	"github.com/txchat/im-util/pkg/net/ws"
	"github.com/txchat/im-util/pkg/user"
	"github.com/txchat/im/api/protocol"
)

const (
	Version0dot2dot0 = "v0.2.0"
)

type ActionInfo struct {
	Err         error           `json:"err"`
	Action      string          `json:"action"`
	UID         string          `json:"uid"`
	ConnID      string          `json:"conn_id"`
	UUID        string          `json:"uuid"`
	From        string          `json:"from"`
	Target      string          `json:"target"`
	ChannelType message.Channel `json:"channel_type"`
	Seq         int32           `json:"seq"`
	Ack         int32           `json:"ack"`
	Mid         string          `json:"mid"`
}

type OnReceiveHandler func(c *net.IMConn, action ActionInfo) error
type OnSendHandler func(c *net.IMConn, action ActionInfo) error

type Device struct {
	u                *user.User
	conn             *net.IMConn
	uuid, deviceName string
	deviceType       auth.Device
	deviceToken      string
	shutdown         chan bool
	isShutdown       int32
	onReceive        OnReceiveHandler
	onSend           OnSendHandler

	httpClient *xhttp.ChatHTTPAPIClient
}

func NewDevice(uuid, deviceName string, deviceType auth.Device, u *user.User) *Device {
	d := &Device{
		u:          u,
		uuid:       uuid,
		deviceName: deviceName,
		deviceType: deviceType,
		shutdown:   make(chan bool),
		isShutdown: 1,
	}
	return d
}

func (d *Device) TurnOn() error {
	if atomic.CompareAndSwapInt32(&d.isShutdown, 1, 0) {
		go d.serve()
	}
	return nil
}

func (d *Device) TurnOff() error {
	if atomic.CompareAndSwapInt32(&d.isShutdown, 0, 1) {
		if d.conn != nil {
			d.conn.Close()
		}
		close(d.shutdown)
	}
	return nil
}

func (d *Device) SetOnReceive(cb OnReceiveHandler) {
	d.onReceive = cb
}

func (d *Device) SetOnSend(cb OnSendHandler) {
	d.onSend = cb
}

func (d *Device) GetUser() *user.User {
	return d.u
}

func (d *Device) DialChatAPI(url url.URL, timeout time.Duration) {
	d.httpClient = xhttp.NewChatHTTPAPIClient(url.String(), timeout, func() *xhttp.AuthenticationMetadata {
		return &xhttp.AuthenticationMetadata{
			Signature:  d.u.Token(),
			UUID:       d.uuid,
			Device:     d.deviceType.String(),
			DeviceName: d.deviceName,
			Version:    Version0dot2dot0,
		}
	})
}

func (d *Device) DialIMServer(appId string, url url.URL, ext []byte) error {
	var authHandler net.AuthHandler
	switch url.Scheme {
	case "ws":
		authHandler = ws.Auth
	case "tcp":
		authHandler = tcp.Auth
	default:
		return fmt.Errorf("unsupported scheme %v", url.Scheme)
	}
	conn, err := net.DialIMAndServe(url.Host, &protocol.AuthBody{
		AppId: appId,
		Token: d.u.Token(),
		Ext:   ext,
	}, 7*time.Second, authHandler)
	if err != nil {
		return err
	}
	d.conn = conn
	return nil
}

func (d *Device) WithDeviceInfo() []byte {
	extData, err := proto.Marshal(&auth.Login{
		Device:      d.deviceType,
		Username:    d.u.GetUsername(),
		DeviceToken: d.deviceToken,
		ConnType:    auth.Login_Connect,
		Uuid:        d.uuid,
		DeviceName:  d.deviceName,
	})
	if err != nil {
		return nil
	}
	return extData
}

func (d *Device) sendMsg(data []byte) (string, error) {
	if d.httpClient == nil {
		return "", errors.New("unset http client")
	}
	resp, err := d.httpClient.SendChatMessage(context.Background(), data)
	if err != nil {
		return "", err
	}
	log.Info().Err(err).Interface("resp", *resp).Msg("sendMsg")
	return xhttp.GetMid(resp)
}

func (d *Device) SendTextMsg(channelType message.Channel, target, text string) error {
	msgType, contentData, err := protoutil.Text(text)
	if err != nil {
		return err
	}
	msg := message.Message{
		ChannelType: channelType,
		Mid:         "",
		Cid:         uuid.New().String(),
		From:        d.u.GetUID(),
		Target:      target,
		MsgType:     msgType,
		Content:     contentData,
		Datetime:    time.Now().UnixMilli(),
		Source:      nil,
		Reference:   nil,
	}

	//发送消息
	data, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}
	mid, err := d.sendMsg(data)
	if d.onSend != nil {
		err = d.onSend(d.conn, ActionInfo{
			Err:         err,
			Action:      "send",
			UID:         d.u.GetUID(),
			ConnID:      d.conn.GetConnId(),
			UUID:        d.uuid,
			From:        d.u.GetUID(),
			Target:      target,
			ChannelType: channelType,
			Seq:         0,
			Ack:         0,
			Mid:         mid,
		})
		if err != nil {
			return err
		}
	}
	return err
}

func (d *Device) serve() {
	for {
		select {
		case <-d.shutdown:
			return
		default:
			revProto := d.conn.Read()
			switch protocol.Op(revProto.GetOp()) {
			case protocol.Op_Message:
				var chatProto chat.Chat
				err := proto.Unmarshal(revProto.GetBody(), &chatProto)
				if err != nil {
					continue
				}
				if chatProto.GetType() == chat.Chat_message {
					var msg message.Message
					err = proto.Unmarshal(chatProto.GetBody(), &msg)
					if err != nil {
						continue
					}
					if msg.GetFrom() != d.u.GetUID() {
						// 记录服务端Push的日志
						// 回调
						if d.onReceive != nil {
							err = d.onReceive(d.conn, ActionInfo{
								Action:      "receive",
								UID:         d.u.GetUID(),
								ConnID:      d.conn.GetConnId(),
								UUID:        d.uuid,
								From:        msg.GetFrom(),
								Target:      msg.GetTarget(),
								ChannelType: msg.GetChannelType(),
								Seq:         revProto.GetSeq(),
								Ack:         revProto.GetAck(),
								Mid:         msg.GetMid(),
							})
							if err != nil {
								continue
							}
						}
					}
				}
				p, err := protoutil.CreateProtoAck(revProto.GetSeq())
				if err != nil {
					continue
				}
				_, err = d.conn.Push(p)
				if err != nil {
					continue
				}
			}
		}
	}
}

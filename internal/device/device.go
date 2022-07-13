package device

import (
	"strconv"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog"
	"github.com/txchat/dtalk/pkg/util"
	protoutil "github.com/txchat/im-util/internal/proto"
	"github.com/txchat/im-util/internal/user"
	"github.com/txchat/im-util/pkg/net"
	"github.com/txchat/im-util/pkg/net/ws"
	comet "github.com/txchat/im/api/comet/grpc"
	"github.com/txchat/im/dtask"
	xproto "github.com/txchat/imparse/proto"
)

type OnReceiveHandler func(c *net.IMConn, proto *comet.Proto) error

type Device struct {
	u                *user.User
	conn             *net.IMConn
	uuid, deviceName string
	deviceType       xproto.Device

	tsk        *dtask.Task
	log        zerolog.Logger
	shutdown   chan bool
	isShutdown int32

	onReceive      OnReceiveHandler
	repushInterval time.Duration
}

func NewDevice(uuid, deviceName string, deviceType xproto.Device, log zerolog.Logger, u *user.User) *Device {
	d := &Device{
		u:              u,
		uuid:           uuid,
		deviceName:     deviceName,
		deviceType:     deviceType,
		shutdown:       make(chan bool),
		isShutdown:     1,
		tsk:            dtask.NewTask(),
		log:            log,
		repushInterval: time.Millisecond * 1500,
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

func (d *Device) SetOnReceive(rev OnReceiveHandler) {
	d.onReceive = rev
}

func (d *Device) GetUser() *user.User {
	return d.u
}

func (d *Device) DialIMServer(appId, server string, ext []byte) error {
	//if atomic.CompareAndSwapInt32(&d.isClosed, 1, 0) {
	//
	//}
	conn, err := net.DialIMAndServe(server, &comet.AuthMsg{
		AppId: appId,
		Token: d.u.Token(),
		Ext:   ext,
	}, 20*time.Second, ws.Auth)
	if err != nil {
		return err
	}
	d.conn = conn
	return nil
}

func (d *Device) DialIMServerWithDeviceInfo(appId, server string) error {
	extData, err := proto.Marshal(&xproto.Login{
		Device:      d.deviceType,
		Username:    d.u.GetUsername(),
		DeviceToken: "",
		ConnType:    xproto.Login_Connect,
		Uuid:        d.uuid,
		DeviceName:  d.deviceName,
	})
	if err != nil {
		return err
	}
	return d.DialIMServer(appId, server, extData)
}

func (d *Device) SendMsg(channelType string, target, text string) {
	//do send
	chType, ok := xproto.Channel_value[channelType]
	if !ok {
		return
	}
	p, err := protoutil.CreateProtoSendMsg(0, d.u.GetUID(), target, xproto.Channel(chType), text)
	if err != nil {
		d.log.Error().Err(err).Msg("CreateProtoSendMsg RangeSend")
		return
	}

	seq, err := d.conn.Push(p)
	if err != nil {
		d.log.Error().Err(err).Msg("Send RangeSend")
		return
	}
	// 发出时间点的日志
	d.log.Info().Str("action", "send").
		Str("user_id", d.u.GetUID()).
		Str("conn_id", d.conn.GetConnId()).
		Int32("seq", seq).
		Str("uuid", d.uuid).
		Str("from", d.u.GetUID()).
		Str("target", target).
		Str("channel_type", channelType).
		Msg("")
	// 添加重传任务
	job, inserted := d.tsk.AddJobRepeat(d.repushInterval, 0, func() {
		_, err := d.conn.RePush(p)
		if err != nil {
			d.log.Error().Err(err).Msg("ReSend RangeSend")
			return
		}
	})
	if !inserted {
		d.log.Error().Err(err).Msg("CreateProtoSendMsg RangeSend")
		return
	}
	d.tsk.Add(util.ToString(seq), job)
}

func (d *Device) serve() {
	for {
		select {
		case <-d.shutdown:
			return
		default:
			// 接收时间点的日志
			revProto := d.conn.Read()
			// do ack
			switch comet.Op(revProto.GetOp()) {
			case comet.Op_ReceiveMsg:
				bizProto, err := protoutil.ConvertBizProto(revProto.GetBody())
				if err != nil {
					d.log.Error().Err(err).Msg("ConvertBizProto Op_ReceiveMsg")
					continue
				}
				if bizProto.GetEventType() == xproto.Proto_common {
					common, err := protoutil.ConvertCommon(bizProto.GetBody())
					if err != nil {
						d.log.Error().Err(err).Msg("ConvertCommon")
						continue
					}
					if common.GetFrom() != d.u.GetUID() {
						// 记录服务端Push的日志
						d.log.Info().Str("action", "receive").
							Str("user_id", d.u.GetUID()).
							Str("conn_id", d.conn.GetConnId()).
							Int64("mid", common.GetMid()).
							Str("from", common.GetFrom()).
							Str("target", common.GetTarget()).
							Str("uuid", d.uuid).
							Msg("")
					}
					// 回调
					if d.onReceive != nil {
						err := d.onReceive(d.conn, revProto)
						if err != nil {
							d.log.Error().Err(err).Msg("onReceive")
							continue
						}
					}
				}

				p, err := protoutil.CreateProtoAck(0, revProto.GetSeq())
				if err != nil {
					d.log.Error().Err(err).Msg("CreateProtoAck Op_ReceiveMsg")
					continue
				}
				_, err = d.conn.Push(p)
				if err != nil {
					d.log.Error().Err(err).Msg("Send Op_ReceiveMsg")
					continue
				}
			case comet.Op_SendMsgReply:
				commAck, err := protoutil.ConvertCommonAck(revProto.GetBody())
				if err != nil {
					d.log.Error().Err(err).Msg("ConvertCommonAck Op_SendMsgReply")
					continue
				}
				// 记录服务端Ack的日志
				d.log.Info().Str("action", "ack").
					Str("user_id", d.u.GetUID()).
					Str("conn_id", d.conn.GetConnId()).
					Int32("ack", revProto.GetAck()).
					Int64("mid", commAck.GetMid()).
					Str("uuid", d.uuid).
					Msg("")
				//从task中删除某一条
				if j := d.tsk.Get(strconv.FormatInt(int64(revProto.GetAck()), 10)); j != nil {
					j.Cancel()
				}
			}
		}
	}
}

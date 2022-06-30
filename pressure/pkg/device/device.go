package device

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/txchat/dtalk/pkg/util"
	"github.com/txchat/im-util/pressure/pkg/msggenerator"
	"github.com/txchat/im-util/pressure/pkg/user"
	comet "github.com/txchat/im/api/comet/grpc"
	"github.com/txchat/im/dtask"
	xproto "github.com/txchat/imparse/proto"
	"strconv"
	"sync/atomic"
	"time"
)

type Device struct {
	u       *user.User
	session *session

	uuid, deviceName string
	deviceType       xproto.Device
	output           chan *xproto.Common
	destroy          chan bool
	isDestroyed      int32
	tsk              *dtask.Task
	log              zerolog.Logger
	repushInterval   time.Duration
}

func NewDevice(uuid, deviceName string, deviceType xproto.Device, log zerolog.Logger, u *user.User) *Device {
	d := &Device{
		u:              u,
		session:        defaultEmpty,
		uuid:           uuid,
		deviceName:     deviceName,
		deviceType:     deviceType,
		output:         make(chan *xproto.Common, 200),
		destroy:        make(chan bool),
		isDestroyed:    0,
		tsk:            dtask.NewTask(),
		log:            log,
		repushInterval: time.Millisecond * 1500,
	}
	go d.serve()
	return d
}

func (d *Device) ConnectIMServer(appId, server string) error {
	log.Info().Str("appId", appId).Str("server", server).Msg("ConnectIMServer")
	return d.u.ConnServerWithDevice(appId, server, &xproto.Login{
		Device:      d.deviceType,
		Username:    d.u.GetUsername(),
		DeviceToken: "",
		ConnType:    xproto.Login_Connect,
		Uuid:        d.uuid,
		DeviceName:  d.deviceName,
	})
}

func (d *Device) DisConnectIMServer() {
	log.Info().Msg("DisConnectIMServer")
	d.u.Close()
}

func (d *Device) ReConnectIMServer(appId, server string) error {
	return d.u.ConnServerWithDevice(appId, server, &xproto.Login{
		Device:      d.deviceType,
		Username:    d.u.GetUsername(),
		DeviceToken: "",
		ConnType:    xproto.Login_Reconnect,
		Uuid:        d.uuid,
		DeviceName:  d.deviceName,
	})
}

func (d *Device) SendMsg(channelType string, target, text string) {
	//do send
	chType, ok := xproto.Channel_value[channelType]
	if !ok {
		return
	}
	log.Info().Interface("user", d.u).Msg("do sendMsg")
	seq := d.u.GenSeq()
	p, err := msggenerator.CreateProtoSendMsg(seq, d.u.GetId(), target, xproto.Channel(chType), text)
	if err != nil {
		d.log.Error().Err(err).Msg("CreateProtoSendMsg RangeSend")
		return
	}
	// 发出时间点的日志
	d.log.Info().Str("device uuid", d.uuid).
		Str("connId", d.u.GetConnId()).
		Str("from", d.u.GetId()).
		Str("target", target).
		Str("channelType", channelType).
		Int32("seq", p.GetSeq()).
		Msg("send msg")
	err = d.u.Send(p)
	if err != nil {
		d.log.Error().Err(err).Msg("Send RangeSend")
		return
	}
	// 添加重传任务
	//push into task pool
	job, inserted := d.tsk.AddJobRepeat(d.repushInterval, 0, func() {
		err := d.u.Send(p)
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

func (d *Device) Destroy() {
	if atomic.CompareAndSwapInt32(&d.isDestroyed, 0, 1) {
		close(d.destroy)
		d.u.Close()
	}
	log.Info().Msg("Device Destroy")
}

func (d *Device) GetIsDestroyed() bool {
	if val := atomic.LoadInt32(&d.isDestroyed); val == 1 {
		return true
	}
	return false
}

func (d *Device) nonblockOutput(p *xproto.Common) {
	select {
	case d.output <- p:
	default:
		d.log.Error().Msg("output device buff full")
	}
}

func (d *Device) GetOutputQueue() chan *xproto.Common {
	return d.output
}

func (d *Device) GetDestroy() chan bool {
	return d.destroy
}

func (d *Device) GetUserId() string {
	return d.u.GetId()
}

func (d *Device) GetSession() *session {
	return d.session
}

func (d *Device) serve() {
	for {
		select {
		case <-d.destroy:
			return
		default:
			// 接收时间点的日志
			revProto := d.u.OnReceive()
			// do ack
			switch comet.Op(revProto.GetOp()) {
			case comet.Op_ReceiveMsg:
				bizProto, err := msggenerator.ConvertBizProto(revProto.GetBody())
				if err != nil {
					d.log.Error().Err(err).Msg("ConvertBizProto Op_ReceiveMsg")
					continue
				}
				if bizProto.GetEventType() == xproto.Proto_common {
					common, err := msggenerator.ConvertCommon(bizProto.GetBody())
					if err != nil {
						d.log.Error().Err(err).Msg("ConvertCommon")
						continue
					}
					// 记录服务端Push的日志
					d.log.Info().Str("device uuid", d.uuid).
						Str("connId", d.u.GetConnId()).
						Str("user_id", d.u.GetId()).
						Int64("mid", common.GetMid()).
						Str("from", common.GetFrom()).
						Str("target", common.GetTarget()).
						Msg("receive message")
					//传下去
					d.nonblockOutput(common)
				}

				p, err := msggenerator.CreateProtoAck(d.u.GenSeq(), revProto.GetSeq())
				if err != nil {
					d.log.Error().Err(err).Msg("CreateProtoAck Op_ReceiveMsg")
					continue
				}
				err = d.u.Send(p)
				if err != nil {
					d.log.Error().Err(err).Msg("Send Op_ReceiveMsg")
					continue
				}
			case comet.Op_SendMsgReply:
				commAck, err := msggenerator.ConvertCommonAck(revProto.GetBody())
				if err != nil {
					d.log.Error().Err(err).Msg("ConvertCommonAck Op_SendMsgReply")
					continue
				}
				// 记录服务端Ack的日志
				d.log.Info().Str("device uuid", d.uuid).
					Str("connId", d.u.GetConnId()).
					Str("user_id", d.u.GetId()).
					Int32("ack", revProto.GetAck()).
					Int64("mid", commAck.GetMid()).
					Msg("receive msg reply")
				//从task中删除某一条
				if j := d.tsk.Get(strconv.FormatInt(int64(revProto.GetAck()), 10)); j != nil {
					j.Cancel()
				}
			}
		}
	}
}

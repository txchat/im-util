package msggenerator

import (
	"github.com/rs/zerolog"
	"github.com/txchat/dtalk/pkg/util"
	"github.com/txchat/im-util/client-pressurev2/pkg/user"
	comet "github.com/txchat/im/api/comet/grpc"
	"github.com/txchat/im/dtask"
	xproto "github.com/txchat/imparse/proto"
	"math/rand"
	"strconv"
	"time"
)

type MsgGenerator struct {
	users     []*user.User
	sendClose chan bool
	ackClose  chan bool
	tsk       *dtask.Task
	log       zerolog.Logger
}

func NewMsgGenerator(users []*user.User, log zerolog.Logger) *MsgGenerator {
	if len(users) < 2 {
		panic("system users less than 2")
	}
	return &MsgGenerator{
		users:     users,
		sendClose: make(chan bool),
		ackClose:  make(chan bool),
		tsk:       dtask.NewTask(),

		log: log,
	}
}

func (m *MsgGenerator) randomTarget(userClient *user.User) string {
	//将时间戳设置成种子数
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(m.users) - 1)
	if m.users[index].GetId() == userClient.GetId() {
		if index >= len(m.users)-1 {
			return m.users[0].GetId()
		}
		return m.users[index+1].GetId()
	}
	return m.users[index].GetId()
}

func (m *MsgGenerator) RangeSend(userClient *user.User, rate time.Duration, sysLog zerolog.Logger) {
	ticker := time.NewTicker(rate)

	for {
		select {
		case <-m.sendClose:
			return
		case <-ticker.C:
			//do send
			seq := userClient.GenSeq()
			p, err := CreateProtoSendMsg(seq, userClient.GetId(), m.randomTarget(userClient), xproto.Channel_ToUser, "1")
			if err != nil {
				sysLog.Error().Err(err).Msg("CreateProtoSendMsg RangeSend")
				continue
			}
			// 发出时间点的日志
			m.log.Info().Str("action", "send").
				Str("user_id", userClient.GetId()).
				Str("conn_id", userClient.GetConnId()).
				Int32("seq", p.GetSeq()).
				Msg("")
			err = userClient.Send(p)
			if err != nil {
				sysLog.Error().Err(err).Msg("Send RangeSend")
				continue
			}
			// 添加重传任务
			//push into task pool
			job, inserted := m.tsk.AddJobRepeat(time.Millisecond*1500, 0, func() {
				err := userClient.Send(p)
				if err != nil {
					sysLog.Error().Err(err).Msg("ReSend RangeSend")
					return
				}
			})
			if !inserted {
				sysLog.Error().Err(err).Msg("CreateProtoSendMsg RangeSend")
				continue
			}
			m.tsk.Add(util.ToString(seq), job)
		}
	}
}

func (m *MsgGenerator) HandleAck(userClient *user.User, sysLog zerolog.Logger) {
	for {
		select {
		case <-m.ackClose:
			return
		default:
			// 接收时间点的日志
			revProto := userClient.OnReceive()
			// do ack
			switch comet.Op(revProto.GetOp()) {
			case comet.Op_ReceiveMsg:
				bizProto, err := ConvertBizProto(revProto.GetBody())
				if err != nil {
					sysLog.Error().Err(err).Msg("ConvertBizProto Op_ReceiveMsg")
					continue
				}
				if bizProto.GetEventType() == xproto.Proto_common {
					common, err := ConvertCommon(bizProto.GetBody())
					if err == nil && common.GetFrom() != userClient.GetId() {
						// 记录服务端Push的日志
						m.log.Info().Str("action", "receive").
							Str("user_id", userClient.GetId()).
							Str("conn_id", userClient.GetConnId()).
							Int64("mid", common.GetMid()).
							Str("from", common.GetFrom()).
							Str("target", common.GetTarget()).
							Msg("")
					}
				}
				p, err := CreateProtoAck(userClient.GenSeq(), revProto.GetSeq())
				if err != nil {
					sysLog.Error().Err(err).Msg("CreateProtoAck Op_ReceiveMsg")
					continue
				}
				err = userClient.Send(p)
				if err != nil {
					sysLog.Error().Err(err).Msg("Send Op_ReceiveMsg")
					continue
				}
			case comet.Op_SendMsgReply:
				commAck, err := ConvertCommonAck(revProto.GetBody())
				if err != nil {
					sysLog.Error().Err(err).Msg("ConvertCommonAck Op_SendMsgReply")
					continue
				}
				// 记录服务端Ack的日志
				m.log.Info().Str("action", "ack").
					Str("user_id", userClient.GetId()).
					Str("conn_id", userClient.GetConnId()).
					Int32("ack", revProto.GetAck()).
					Int64("mid", commAck.GetMid()).
					Msg("")
				//从task中删除某一条
				if j := m.tsk.Get(strconv.FormatInt(int64(revProto.GetAck()), 10)); j != nil {
					j.Cancel()
				}
			}
		}
	}
}

func (m *MsgGenerator) StopSend() {
	close(m.sendClose)
}

func (m *MsgGenerator) StopAck() {
	close(m.ackClose)
}

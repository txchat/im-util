package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/types"
	"github.com/33cn/chain33/wallet/bipwallet"
	"github.com/rs/zerolog"
	protoutil "github.com/txchat/im-util/internal/proto"
	"github.com/txchat/im-util/internal/user"
	"github.com/txchat/im-util/pkg/net"
	"github.com/txchat/im-util/pkg/net/ws"
	comet "github.com/txchat/im/api/comet/grpc"
	xproto "github.com/txchat/imparse/proto"
)

var (
	log           zerolog.Logger
	appId, server string
)

func init() {
	log = zerolog.New(os.Stdout).With().Timestamp().Logger()
	flag.StringVar(&appId, "appId", "dtalk", "")
	flag.StringVar(&server, "server", "127.0.0.1:3102", "")
}

func main() {
	flag.Parse()
	seq := int32(0)
	mid := int64(0)
	closer := make(chan bool)
	log.Info().Str("appId", appId).Str("server", server).Msg("start")

	a, err := newUser()
	if err != nil {
		log.Error().Err(err).Msg("create user a failed")
		os.Exit(1)
	}
	connA, err := dial(a, func(conn *net.IMConn) {
		for {
			select {
			case <-closer:
				return
			default:
				revProto := conn.Read()
				// do ack
				switch comet.Op(revProto.GetOp()) {
				case comet.Op_ReceiveMsg:
					bizProto, err := protoutil.ConvertBizProto(revProto.GetBody())
					if err != nil {
						log.Error().Err(err).Msg("ConvertBizProto Op_ReceiveMsg")
						continue
					}
					if bizProto.GetEventType() == xproto.Proto_common {
						common, err := protoutil.ConvertCommon(bizProto.GetBody())
						if err == nil && common.GetFrom() != a.GetUID() {
							//// 记录服务端Push的日志
							//log.Info().Str("action", "receive").
							//	Str("user_id", a.GetId()).
							//	Str("conn_id", a.GetConnId()).
							//	Int64("mid", common.GetMid()).
							//	Str("from", common.GetFrom()).
							//	Str("target", common.GetTarget()).
							//	Msg("")
						}
					}
					p, err := protoutil.CreateProtoAck(0, revProto.GetSeq())
					if err != nil {
						log.Error().Err(err).Msg("CreateProtoAck Op_ReceiveMsg")
						continue
					}
					_, err = conn.Push(p)
					if err != nil {
						log.Error().Err(err).Msg("Send Op_ReceiveMsg")
						continue
					}
				case comet.Op_SendMsgReply:
					commAck, err := protoutil.ConvertCommonAck(revProto.GetBody())
					if err != nil {
						log.Error().Err(err).Msg("ConvertCommonAck Op_SendMsgReply")
						continue
					}
					// 记录服务端Ack的日志
					if seq == revProto.GetAck() {
						mid = commAck.GetMid()
						log.Info().Msg("接收发送响应成功，步骤2达成")
					}
				}
			}
		}
	})
	if err != nil {
		log.Error().Err(err).Msg("user a conn server failed")
		os.Exit(1)
	}

	b, err := newUser()
	if err != nil {
		log.Error().Err(err).Msg("create user b failed")
		os.Exit(1)
	}
	_, err = dial(b, func(conn *net.IMConn) {
		for {
			select {
			case <-closer:
				return
			default:
				// 接收时间点的日志
				revProto := conn.Read()
				// do ack
				switch comet.Op(revProto.GetOp()) {
				case comet.Op_ReceiveMsg:
					bizProto, err := protoutil.ConvertBizProto(revProto.GetBody())
					if err != nil {
						log.Error().Err(err).Msg("ConvertBizProto Op_ReceiveMsg")
						continue
					}
					if bizProto.GetEventType() == xproto.Proto_common {
						common, err := protoutil.ConvertCommon(bizProto.GetBody())
						if err == nil && common.GetFrom() != b.GetUID() {
							if mid == common.GetMid() {
								// 记录服务端Push的日志
								log.Info().Msg("对端接收成功，步骤3达成")
							}
						}
					}
					p, err := protoutil.CreateProtoAck(0, revProto.GetSeq())
					if err != nil {
						log.Error().Err(err).Msg("CreateProtoAck Op_ReceiveMsg")
						continue
					}
					_, err = conn.Push(p)
					if err != nil {
						log.Error().Err(err).Msg("Send Op_ReceiveMsg")
						continue
					}
				case comet.Op_SendMsgReply:
				}
			}
		}
	})
	if err != nil {
		log.Error().Err(err).Msg("user b conn server failed")
		os.Exit(1)
	}

	log.Info().Msg(`
	待确认3个步骤：
步骤1：a成功发送消息
步骤2：a接收发送响应
步骤3：b接收到a发送的数据
`)

	seq, err = SendMsg(connA, xproto.Channel_name[int32(xproto.Channel_ToUser)], a.GetUID(), b.GetUID(), "1")
	// create proto
	if err != nil {
		log.Error().Err(err).Msg("user a send proto failed")
		os.Exit(1)
	}

	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			close(closer)
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}

func newUser() (*user.User, error) {
	//创建助记词
	mne, err := bipwallet.NewMnemonicString(1, 160)
	if err != nil {
		return nil, err
	}
	//创建钱包
	wallet, err := bipwallet.NewWalletFromMnemonic(bipwallet.TypeBty, uint32(types.SECP256K1), mne)
	if err != nil {
		return nil, err
	}
	private, public, err := wallet.NewKeyPair(0)
	if err != nil {
		return nil, err
	}
	addr := address.PubKeyToAddr(0, public)
	return user.NewUser(addr, private, public), nil
}

func dial(u *user.User, cb func(conn *net.IMConn)) (*net.IMConn, error) {
	conn, err := net.DialIMAndServe(server, &comet.AuthMsg{
		AppId: appId,
		Token: u.Token(),
		Ext:   nil,
	}, 20*time.Second, ws.Auth)
	if err != nil {
		return nil, err
	}
	go cb(conn)
	return conn, nil
}

func SendMsg(conn *net.IMConn, channelType string, from, target, text string) (int32, error) {
	//do send
	chType, ok := xproto.Channel_value[channelType]
	if !ok {
		return 0, fmt.Errorf("undefined channel type %v", chType)
	}
	p, err := protoutil.CreateProtoSendMsg(0, from, target, xproto.Channel(chType), text)
	if err != nil {
		log.Error().Err(err).Msg("CreateProtoSendMsg RangeSend")
		return 0, err
	}

	seq, err := conn.Push(p)
	if err != nil {
		log.Error().Err(err).Msg("Send RangeSend")
		return 0, err
	}
	log.Info().Msg("发送成功，步骤1达成")
	return seq, nil
}

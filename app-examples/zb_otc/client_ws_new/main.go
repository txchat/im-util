package main

// ./client {appId} {token} {ws-server-addr}
// ./client echo f3dc8ccd localhost:3102

import (
	"flag"
	"fmt"
	"github.com/txchat/im-util/lib/logger"
	"github.com/txchat/im-util/lib/ws"
	comet "github.com/txchat/im/api/comet/grpc"
	"runtime"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/txchat/im-util/lib"
	bizProto "github.com/txchat/imparse/proto"
)

var log = logger.Log

const (
	ping = 2 * time.Second
)

var (
	debug   bool
	appId   string
	token   string
	address string
	num     int
)

func init() {
	flag.StringVar(&appId, "appid", "", "sets app id")
	flag.StringVar(&token, "token", "", "sets token")
	flag.StringVar(&address, "address", "", "sets address")
	flag.BoolVar(&debug, "debug", false, "sets log level to debug")
	flag.IntVar(&num, "n", 1, "user number")
}

// uid = 18
// send -appid=zb_otc -token="Bearer 5a40325c991f254441bdda66df1c1241516e7fe1" -address=47.242.199.130:3102 -debug=true
// send -appid=zb_otc -token="Bearer 5a40325c991f254441bdda66df1c1241516e7fe1" -address=172.16.101.126:3102 -debug=true
// send -appid=zb_otc -token="Bearer 5a40325c991f254441bdda66df1c1241516e7fe1" -address=127.0.0.1:3102 -debug=true
// uid = 3
// receive -appid=zb_otc -token="Bearer 68cfae25cd487734ab20236b8596187277a7c8fd" -address=47.242.199.130:3102 -debug=true
// receive -appid=zb_otc -token="Bearer 68cfae25cd487734ab20236b8596187277a7c8fd" -address=172.16.101.126:3102 -debug=true
// receive -appid=zb_otc -token="Bearer 68cfae25cd487734ab20236b8596187277a7c8fd" -address=127.0.0.1:3102 -debug=true

func initLog(debug bool) {
	logger.Init(debug)
	lib.InitLog(debug, "")
	log = logger.Log
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	initLog(debug)

	fmt.Println(appId, token, address, num)
	for i := 0; i < num; i++ {
		go client(appId, token, address)
	}
	var exit chan bool
	<-exit
}

func client(appId, token, server string) {
	cli, err := lib.NewClient(appId, token, server, 5*time.Second, ws.Auth)
	if err != nil {
		panic(err)
	}
	cli.SetBiz(new(biz))
	cli.Serve()
	Send(cli)
}

type biz struct {
}

func (b *biz) Receive(c *lib.Client, p *comet.Proto) error {

	receivedProto := &bizProto.Proto{}
	if err := proto.Unmarshal(p.GetBody(), receivedProto); err != nil {
		log.Error().Err(err).Msg("receivedProto Unmarshal failed")
		return err
	}

	receivedCommonMsg := &bizProto.CommonMsg{}
	if err := proto.Unmarshal(receivedProto.GetBody(), receivedCommonMsg); err != nil {
		log.Error().Err(err).Msg("receivedCommonMsg Unmarshal failed")
		return err
	}
	log.Info().Str("op", comet.Op_name[p.Op]).
		Str("msg", string(receivedCommonMsg.GetMsg())).
		Str("from", receivedCommonMsg.GetFrom()).
		Str("target", receivedCommonMsg.GetTarget()).
		Msg("Receive")

	if p.Op == int32(comet.Op_ReceiveMsg) {
		pp := new(comet.Proto)
		pp.Op = int32(comet.Op_ReceiveMsgReply)
		pp.Seq = c.IncSeq()
		pp.Ack = p.GetSeq()
		body, err := proto.Marshal(TextReceiveMsgReply(""))
		if err != nil {
			return nil
		}

		pp.Body = body
		if err = c.Push(pp); err != nil {
			return nil
		}

	}

	return nil
}

func Send(c *lib.Client) {
	//write message
	go func() {
		pp := new(comet.Proto)
		for {
			time.Sleep(ping)
			pp.Op = int32(comet.Op_SendMsg)
			pp.Seq = c.IncSeq()
			pp.Ack = 0
			baseMsg := TextMsg(time.Now().Format(time.RFC3339))
			body, err := proto.Marshal(baseMsg)
			if err != nil {
				return
			}

			pp.Body = body
			if err = c.Push(pp); err != nil {
				return
			}
		}
	}()
}

//
func TextMsg(msg string) *bizProto.Proto {
	m := bizProto.TextMsg{
		Content: msg,
	}
	log.Info().Str("msg", msg).
		Str("from", "953").
		Str("target", "20210419184685527C1DB08A").
		Msg("send")
	data, err := proto.Marshal(&m)
	if err != nil {
		panic(err)
	}
	body, err := proto.Marshal(&bizProto.CommonMsg{
		ChannelType: int32(bizProto.ToUser),
		LogId:       0,
		MsgId:       uuid.New().String(),
		From:        "953",
		Target:      "20210419184685527C1DB08A_Bearer 5a40325c991f254441bdda66df1c1241516e7fe1",
		MsgType:     int32(bizProto.MsgType_Text),
		Msg:         data,
		Datetime:    0,
	})
	if err != nil {
		panic(err)
	}
	return &bizProto.Proto{
		EventType: bizProto.EventType_commonMsg,
		Body:      body,
	}
}

func TextReceiveMsgReply(msg string) *bizProto.Proto {
	m := bizProto.TextMsg{
		Content: msg,
	}
	log.Info().Str("msg", msg).
		Str("from", "40").
		Str("target", "20210419184685527C1DB08A").
		Msg("send")
	data, err := proto.Marshal(&m)
	if err != nil {
		panic(err)
	}
	body, err := proto.Marshal(&bizProto.CommonMsg{
		ChannelType: int32(bizProto.ToUser),
		LogId:       0,
		MsgId:       uuid.New().String(),
		From:        "40",
		Target:      "20210419184685527C1DB08A_Bearer 5a40325c991f254441bdda66df1c1241516e7fe1",
		MsgType:     int32(bizProto.MsgType_Text),
		Msg:         data,
		Datetime:    0,
	})
	if err != nil {
		panic(err)
	}
	return &bizProto.Proto{
		EventType: bizProto.EventType_commonMsgAck,
		Body:      body,
	}
}

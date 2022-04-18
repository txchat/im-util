package main

// ./client {appId} {token} {ws-server-addr}
// ./client echo f3dc8ccd localhost:3102

import (
	"flag"
	"fmt"
	"github.com/txchat/im-util/lib/ws"
	comet "github.com/txchat/im/api/comet/grpc"
	"os"
	"runtime"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/txchat/im-util/lib"
	bizProto "github.com/txchat/imparse/proto"
)

const (
	ping = 2 * time.Second
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	go client(os.Args[1], os.Args[2], os.Args[3])
	var exit chan bool
	<-exit
}

func client(appId, token, server string) {
	cli, err := lib.NewClient(appId, token, server, 20*time.Second, ws.Auth)
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
	return nil
}

func Send(c *lib.Client) {
	//write message
	go func() {
		pp := new(comet.Proto)
		for {
			time.Sleep(ping)
			fmt.Println("push one")
			pp.Op = int32(comet.Op_SendMsg)
			pp.Seq = c.IncSeq()
			pp.Ack = 0
			baseMsg := TextMsg("hello")
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
	data, err := proto.Marshal(&m)
	if err != nil {
		panic(err)
	}
	body, err := proto.Marshal(&bizProto.CommonMsg{
		ChannelType: int32(bizProto.ToUser),
		LogId:       0,
		MsgId:       uuid.New().String(),
		From:        "1FKxgaEh5fuSm7a35BfUnKYAmradowpiTR",
		Target:      "1AsPsahP7FvpR7F2de1LhSB4SU5ShqZ7eu",
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

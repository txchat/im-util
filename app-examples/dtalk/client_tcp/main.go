package main

// ./client {appId} {token} {ws-server-addr}
// ./client echo f3dc8ccd localhost:3102

import (
	"flag"
	"os"
	"runtime"
	"time"

	protoutil "github.com/txchat/im-util/internal/proto"
	"github.com/txchat/im-util/pkg/net"
	"github.com/txchat/im-util/pkg/net/tcp"
	comet "github.com/txchat/im/api/comet/grpc"
	xproto "github.com/txchat/imparse/proto"
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
	cli, err := net.DialIMAndServe(server, &comet.AuthMsg{
		AppId: appId,
		Token: token,
		Ext:   nil,
	}, 5*time.Second, tcp.Auth)
	if err != nil {
		panic(err)
	}
	Send(cli)
}

func Send(c *net.IMConn) {
	//write message
	go func() {
		pp, err := protoutil.CreateProtoSendMsg(0, "1FKxgaEh5fuSm7a35BfUnKYAmradowpiTR", "1AsPsahP7FvpR7F2de1LhSB4SU5ShqZ7eu", xproto.Channel_ToUser, "hello")
		if err != nil {
			return
		}
		for {
			time.Sleep(ping)
			if _, err = c.Push(pp); err != nil {
				return
			}
		}
	}()
}

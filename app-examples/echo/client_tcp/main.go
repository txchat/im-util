package main

// ./client {appId} {token} {tcp-server-addr}
// ./client echo f3dc8ccd localhost:3101

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/Terry-Mao/goim/pkg/bufio"
	"github.com/golang/protobuf/proto"
	"github.com/txchat/im-util/app-examples/echo/types"
	"github.com/txchat/im/api/protocol"
)

const (
	heart = 5 * time.Second
	ping  = 2 * time.Second
)

var (
	lck sync.RWMutex
	seq int32
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	go client(os.Args[1], os.Args[2], os.Args[3])
	var exit chan bool
	<-exit
}

func client(appId, token, server string) {
	quit := make(chan bool, 1)
	defer func() {
		close(quit)
	}()

	// connnect to server
	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Printf("net.Dial(%s) error(%v)", server, err)
		return
	}
	wr := bufio.NewWriter(conn)
	rd := bufio.NewReader(conn)
	authMsg := &protocol.AuthBody{
		AppId: appId,
		Token: token,
	}
	p := new(protocol.Proto)
	p.Ver = 1
	p.Op = int32(protocol.Op_Auth)
	p.Seq = seq
	p.Body, _ = proto.Marshal(authMsg)

	//auth
	if err = WriteProto(wr, p); err != nil {
		log.Printf("auth error(%v)", err)
		return
	}
	if err = p.ReadTCP(rd); err != nil {
		log.Printf("tcpReadProto() error(%v)", err)
		return
	}
	uid := fmt.Sprintf("%s:%s", appId, token)
	log.Printf("%s auth ok", uid)
	seq++

	// writer heartbeat
	go func() {
		hb := new(protocol.Proto)
		for {
			// heartbeat
			hb.Op = int32(protocol.Op_Heartbeat)
			hb.Seq = GetSeq()
			hb.Body = nil
			if err = WriteProto(wr, hb); err != nil {
				log.Printf("uid:%s hb error(%v)", uid, err)
				return
			}
			log.Printf("uid:%s Write heartbeat", uid)
			time.Sleep(heart)
			IncSeq()
			select {
			case <-quit:
				return
			default:
			}
		}
	}()

	// write ping
	go func() {
		pp := new(protocol.Proto)
		for {
			pp.Op = int32(protocol.Op_Message)
			pp.Seq = GetSeq()
			msg := fmt.Sprintf("ping:%v", pp.Seq)
			echo := &types.EchoMsg{
				Value: &types.EchoMsg_Ping{Ping: &types.Ping{Msg: msg}},
				Ty:    int32(types.EchoOp_PingAction),
			}
			body, err := proto.Marshal(echo)
			if err != nil {
				log.Printf("marshal error(%v)", err)
				return
			}

			pp.Body = body
			if err = WriteProto(wr, pp); err != nil {
				log.Printf("uid:%s send ping error(%v)", uid, err)
				return
			}
			xxx, _ := json.Marshal(pp)
			log.Printf("uid:%s Write msg:[%v], Ptoto:[%v]", uid, echo, string(xxx))
			time.Sleep(ping)
			IncSeq()
			select {
			case <-quit:
				return
			default:
			}
		}
	}()

	// reader
	for {
		if err = p.ReadTCP(rd); err != nil {
			log.Printf("uid:%s tcpReadProto() error(%v)", uid, err)
			quit <- true
			return
		}

		xxx, _ := json.Marshal(p)
		log.Printf("server resp p:[%v]", string(xxx))

		if p.Op == int32(protocol.Op_AuthReply) {
			log.Printf("uid:%s auth success", uid)
		} else if p.Op == int32(protocol.Op_HeartbeatReply) {
			log.Printf("uid:%s recv heartbeat reply", uid)
			if err = conn.SetReadDeadline(time.Now().Add(heart + 60*time.Second)); err != nil {
				log.Printf("conn.SetReadDeadline() error(%v)", err)
				quit <- true
				return
			}
		}
	}
}

func IncSeq() {
	lck.Lock()
	seq++
	lck.Unlock()
}

func GetSeq() int32 {
	lck.RLock()
	defer lck.RUnlock()
	return seq
}

func WriteProto(wr *bufio.Writer, p *protocol.Proto) error {
	lck.Lock()
	defer lck.Unlock()

	if err := p.WriteTCP(wr); err != nil {
		log.Printf("tcpWriteProto() error(%v)", err)
		return err
	}
	return wr.Flush()
}

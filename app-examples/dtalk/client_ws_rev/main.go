package main

// ./client {appId} {token} {ws-server-addr}
// ./client echo f3dc8ccd localhost:3102

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	comet "github.com/txchat/im/api/comet/grpc"
	bizProto "github.com/txchat/imparse/proto"
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
	wsUrl := "ws://" + os.Args[3] + "/sub"
	fmt.Println("wsUrl", wsUrl)
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		panic(err)
	}

	authMsg := &comet.AuthMsg{
		AppId: appId,
		Token: token,
	}
	p := new(comet.Proto)
	p.Ver = 1
	p.Op = int32(comet.Op_Auth)
	p.Seq = seq
	p.Body, _ = proto.Marshal(authMsg)

	//auth
	if err = wsWriteProto(conn, p); err != nil {
		log.Printf("wsWriteProto() error(%v)", err)
		return
	}
	if err = wsReadProto(conn, p); err != nil {
		log.Printf("tcpReadProto() error(%v)", err)
		return
	}
	uid := fmt.Sprintf("%s:%s", appId, token)
	log.Printf("%s auth ok", uid)
	seq++

	// writer heartbeat
	go func() {
		hbProto := new(comet.Proto)
		for {
			// heartbeat
			hbProto.Op = int32(comet.Op_Heartbeat)
			hbProto.Seq = GetSeq()
			hbProto.Body = nil
			if err = wsWriteProto(conn, hbProto); err != nil {
				log.Printf("uid:%s tcpWriteProto() error(%v)", uid, err)
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
	//go func() {
	//	pp := new(comet.Proto)
	//	for {
	//		pp.Op = int32(comet.Op_SendMsg)
	//		pp.Seq = GetSeq()
	//		//msg := fmt.Sprintf("ping:%v", pp.Seq)
	//		baseMsg := TextMsg("你好")
	//		body, err := proto.Marshal(baseMsg)
	//		if err != nil {
	//			log.Printf("marshal error(%v)", err)
	//			return
	//		}
	//
	//		pp.Body = body
	//		if err = wsWriteProto(conn, pp); err != nil {
	//			log.Printf("uid:%s tcpWriteProto() error(%v)", uid, err)
	//			return
	//		}
	//		xxx, _ := json.Marshal(pp)
	//		log.Printf("uid:%s Write msg:[%v], Ptoto:[%v]", uid, baseMsg, string(xxx))
	//		time.Sleep(ping)
	//		IncSeq()
	//		select {
	//		case <-quit:
	//			return
	//		default:
	//		}
	//	}
	//}()

	// reader
	for {
		if err = wsReadProto(conn, p); err != nil {
			log.Printf("uid:%s tcpReadProto() error(%v)", uid, err)
			quit <- true
			return
		}

		xxx, _ := json.Marshal(p)
		log.Printf("server resp p:[%v]", string(xxx))

		if p.Op == int32(comet.Op_AuthReply) {
			log.Printf("uid:%s auth success", uid)
		} else if p.Op == int32(comet.Op_HeartbeatReply) {
			log.Printf("uid:%s recv heartbeat reply", uid)
			if err = conn.SetReadDeadline(time.Now().Add(heart + 60*time.Second)); err != nil {
				log.Printf("conn.SetReadDeadline() error(%v)", err)
				quit <- true
				return
			}
		} else {
			// biz msg resp
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

func wsWriteProto(conn *websocket.Conn, proto *comet.Proto) (err error) {
	lck.Lock()
	defer lck.Unlock()
	return proto.WriteWebsocket2(conn)

	//wc, err := conn.NextWriter(websocket.BinaryMessage)
	//if err != nil {
	//	panic(err)
	//}
	//wr := bufio.NewWriter(wc)
	//
	//// write
	//if err = binary.Write(wr, binary.BigEndian, uint32(rawHeaderLen)+uint32(len(proto.Body))); err != nil {
	//	return
	//}
	//if err = binary.Write(wr, binary.BigEndian, rawHeaderLen); err != nil {
	//	return
	//}
	//if err = binary.Write(wr, binary.BigEndian, int16(proto.Ver)); err != nil {
	//	return
	//}
	//if err = binary.Write(wr, binary.BigEndian, proto.Op); err != nil {
	//	return
	//}
	//if err = binary.Write(wr, binary.BigEndian, proto.Seq); err != nil {
	//	return
	//}
	//if proto.Body != nil {
	//	if err = binary.Write(wr, binary.BigEndian, proto.Body); err != nil {
	//		return
	//	}
	//}
	//err = wr.Flush()
	//wc.Close()
	//return
}

func wsReadProto(conn *websocket.Conn, proto *comet.Proto) (err error) {
	return proto.ReadWebsocket2(conn)
	//var (
	//	packLen   int32
	//	headerLen int16
	//	ver       int16
	//)
	//
	//_, rc, err := conn.NextReader()
	//if err != nil {
	//	log.Printf("NextReader error(%v) %s", err)
	//	return err
	//}
	//rd := bufio.NewReader(rc)
	//
	//// read
	//if err = binary.Read(rd, binary.BigEndian, &packLen); err != nil {
	//	return
	//}
	//if err = binary.Read(rd, binary.BigEndian, &headerLen); err != nil {
	//	return
	//}
	//if err = binary.Read(rd, binary.BigEndian, &ver); err != nil {
	//	proto.Ver = int32(ver)
	//	return
	//}
	//if err = binary.Read(rd, binary.BigEndian, &proto.Op); err != nil {
	//	return
	//}
	//if err = binary.Read(rd, binary.BigEndian, &proto.Seq); err != nil {
	//	return
	//}
	//var (
	//	n, t    int
	//	bodyLen = int(packLen - int32(headerLen))
	//)
	//if bodyLen > 0 {
	//	proto.Body = make([]byte, bodyLen)
	//	for {
	//		if t, err = rd.Read(proto.Body[n:]); err != nil {
	//			return
	//		}
	//		if n += t; n == bodyLen {
	//			break
	//		}
	//	}
	//} else {
	//	proto.Body = nil
	//}
	//return
}

//
func TextMsg(msg string) *bizProto.Proto {
	m := map[string]string{
		"content": msg,
	}
	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return &bizProto.Proto{
		EventType:   0,
		ChannelType: 0,
		LogId:       0,
		From:        "1FKxgaEh5fuSm7a35BfUnKYAmradowpiTR",
		Target:      "1AsPsahP7FvpR7F2de1LhSB4SU5ShqZ7eu",
		MsgType:     1,
		Msg:         data,
	}
}

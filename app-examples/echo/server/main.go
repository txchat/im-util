package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	cluster "github.com/bsm/sarama-cluster"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/txchat/im-util/app-examples/echo/types"
	comet "github.com/txchat/im/api/comet/grpc"
	logic "github.com/txchat/im/api/logic/grpc"
	"github.com/txchat/im/common"
)

var (
	appId      string
	logicAddr  string
	mqAddr     string
	listenAddr string
)

type Job struct {
	consumer    *cluster.Consumer
	logicClient logic.LogicClient
}

func init() {
	flag.StringVar(&appId, "app-server name", "echo", "msg queue addr")
	flag.StringVar(&logicAddr, "logic", "127.0.0.1:3119", "logic rpc server addr")
	flag.StringVar(&mqAddr, "mq", "172.16.101.107:9092", "msg queue addr")
	flag.StringVar(&listenAddr, "app-server listen addr", ":9999", "msg queue addr")
}

// New new a push job.
func New() *Job {
	j := &Job{
		consumer:    newKafkaSub(),
		logicClient: newLogicClient(logicAddr),
	}

	return j
}

func newKafkaSub() *cluster.Consumer {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	topic := fmt.Sprintf("goim-%s-topic", appId)
	group := fmt.Sprintf("goim-%s-group", appId)
	consumer, err := cluster.NewConsumer([]string{mqAddr}, group, []string{topic}, config)
	if err != nil {
		panic(err)
	}
	return consumer
}

func newLogicClient(addr string) logic.LogicClient {
	conn, err := common.NewGRPCConn(addr, time.Second)
	if err != nil {
		panic(err)
	}
	return logic.NewLogicClient(conn)
}

// Consume messages, watch signals
func (j *Job) Consume() {
	for {
		select {
		case err := <-j.consumer.Errors():
			log.Printf("consumer error(%v)", err)
		case n := <-j.consumer.Notifications():
			log.Printf("consumer rebalanced(%v)", n)
		case msg, ok := <-j.consumer.Messages():
			if !ok {
				return
			}
			j.consumer.MarkOffset(msg, "")
			bizMsg := new(logic.BizMsg)
			if err := proto.Unmarshal(msg.Value, bizMsg); err != nil {
				log.Printf("proto.Unmarshal(%v) error(%v)", msg, err)
				continue
			}
			//log.Printf("consume: %s/%d/%d\t%s\t%+v", msg.Topic, msg.Partition, msg.Offset, msg.Key, bizMsg)
			j.processMsg(bizMsg)
		}
	}
}

func (j *Job) processMsg(m *logic.BizMsg) {
	if m.AppId != appId {
		return
	}

	var p comet.Proto
	err := proto.Unmarshal(m.Msg, &p)
	if err != nil {
		log.Printf("unmarshal proto error %v", err)
		return
	}

	if m.Op == int32(comet.Op_Auth) {
		log.Printf("user %v login  with key %s", m.FromId, m.Key)
	} else if m.Op == int32(comet.Op_Disconnect) {
		log.Printf("user %v logout with key %s", m.FromId, m.Key)
	} else {
		xxx, _ := json.Marshal(p)
		log.Printf("recv Proto:[%v]", string(xxx))

		var echo types.EchoMsg
		err = proto.Unmarshal(p.Body, &echo)
		if err != nil {
			log.Printf("unmarshal echo error(%v) ", err)
			return
		}
		xxx, _ = json.Marshal(echo)
		log.Printf("client echo:[%v]", string(xxx))

		p.Body, err = makeResp(&echo)
		if err != nil {
			log.Printf("makeResp error(%v) ", err)
			return
		}

		xxx, _ = json.Marshal(p)
		log.Printf("resp Proto:[%v]", string(xxx))
		bytes, _ := proto.Marshal(&p)
		keysMsg := &logic.KeysMsg{
			AppId:  appId,
			ToKeys: []string{m.Key},
			Msg:    bytes,
		}

		_, err = j.logicClient.PushByKeys(context.Background(), keysMsg)
		if err != nil {
			log.Printf("PushByKeys %s, %s, (%v) ", m.FromId, m.Key, err)
		}
	}
}

func makeResp(e *types.EchoMsg) ([]byte, error) {
	var ee types.EchoMsg
	ee.Ty = int32(types.EchoOp_PangAction)
	ee.Value = &types.EchoMsg_Pang{&types.Pang{Msg: fmt.Sprintf("pang for %s", e.Value.(*types.EchoMsg_Ping).Ping.Msg)}}
	return proto.Marshal(&ee)
}

func main() {
	flag.Parse()

	e := gin.Default()
	e.POST("/login", login)
	go e.Run(listenAddr)

	j := New()
	go j.Consume()

	select {}
}

func login(c *gin.Context) {
	head := c.GetHeader("Authorization")
	ss := strings.Split(head, " ")
	if len(ss) < 2 || ss[0] != "Bearer" {
		ret := map[string]interface{}{
			"code":  -1,
			"error": map[string]interface{}{"message": "invalid header"},
		}
		c.PureJSON(http.StatusOK, ret)
		return
	}

	id := fmt.Sprintf("%v", time.Now().UnixNano())
	ret := map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{"id": id},
	}
	c.PureJSON(http.StatusOK, ret)
}

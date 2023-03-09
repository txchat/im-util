package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	xkafka "github.com/oofpgDLD/kafka-go"
	"github.com/txchat/im-util/app-examples/echo/types"
	"github.com/txchat/im/api/protocol"
	"github.com/txchat/im/app/logic/logicclient"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

var (
	appId      string
	logicAddr  string
	mqAddr     string
	listenAddr string
)

type Job struct {
	batchConsumer *xkafka.BatchConsumer
	logicClient   logicclient.Logic
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
		logicClient: logicclient.NewLogic(zrpc.MustNewClient(zrpc.RpcClientConf{
			Endpoints: []string{logicAddr},
			Timeout:   2000,
		}, zrpc.WithNonBlock())),
	}

	//new batch consumer
	consumer := xkafka.NewConsumer(xkafka.ConsumerConfig{
		Version:        "",
		Brokers:        []string{mqAddr},
		Group:          fmt.Sprintf("goim-%s-receive", appId),
		Topic:          fmt.Sprintf("goim-%s-receive-echo", appId),
		CacheCapacity:  100,
		ConnectTimeout: 100,
	}, nil)
	logx.Info("dial kafka broker success")
	bc := xkafka.NewBatchConsumer(xkafka.BatchConsumerConf{
		CacheCapacity: 100,
		Consumers:     100,
		Processors:    100,
	}, xkafka.WithHandle(j.process), consumer)
	j.batchConsumer = bc
	return j
}

func (j *Job) process(key string, data []byte) error {
	ctx := context.Background()
	receivedMsg := new(logicclient.ReceivedMessage)
	if err := proto.Unmarshal(data, receivedMsg); err != nil {
		return err
	}
	if receivedMsg.GetAppId() != appId {
		log.Printf("unsupported appID")
		return fmt.Errorf("unsupported appID")
	}

	switch receivedMsg.GetOp() {
	case protocol.Op_Message:
		var echo types.EchoMsg
		err := proto.Unmarshal(receivedMsg.GetBody(), &echo)
		if err != nil {
			log.Printf("echoMsg unmarshal error(%v) ", err)
			return err
		}

		receivedMsg.Body, err = makeResp(&echo)
		if err != nil {
			log.Printf("makeResp error(%v) ", err)
			return err
		}

		bytes, err := proto.Marshal(receivedMsg)
		if err != nil {
			log.Printf("receivedMsg marshal error(%v) ", err)
			return err
		}

		_, err = j.logicClient.PushByKey(ctx, &logicclient.PushByKeyReq{
			AppId: appId,
			ToKey: []string{receivedMsg.Key},
			Op:    receivedMsg.GetOp(),
			Body:  bytes,
		})
		if err != nil {
			log.Printf("echo logic push error(%v) ", err)
			return err
		}
	default:
		return fmt.Errorf("received message operation %v unsupported", receivedMsg.GetOp())
	}
	return nil
}

func makeResp(e *types.EchoMsg) ([]byte, error) {
	var ee types.EchoMsg
	ee.Ty = int32(types.EchoOp_PongAction)
	ee.Value = &types.EchoMsg_Pong{Pong: &types.Pong{Msg: fmt.Sprintf("pong for %s", e.Value.(*types.EchoMsg_Ping).Ping.Msg)}}
	return proto.Marshal(&ee)
}

func main() {
	flag.Parse()

	e := gin.Default()
	e.POST("/login", login)
	go e.Run(listenAddr)

	j := New()

	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			j.batchConsumer.GracefulStop(ctx)
			cancel()
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
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

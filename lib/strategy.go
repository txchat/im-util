package lib

import (
	comet "github.com/txchat/im/api/comet/grpc"
	"time"
)

type Biz interface {
	Receive(c *Client, proto *comet.Proto) error
	//Send() proto.Message
}

type ReaderWriter interface {
	WriteProto(proto *comet.Proto) (err error)
	ReadProto(proto *comet.Proto) (err error)
	SetReadDeadline(t time.Time) error
}

type AuthHandler func(appId, token, server string, ext []byte) (ReaderWriter, error)

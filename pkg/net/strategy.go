package net

import (
	comet "github.com/txchat/im/api/comet/grpc"
	"time"
)

type ReaderWriterCloser interface {
	WriteProto(proto *comet.Proto) (err error)
	ReadProto(proto *comet.Proto) (err error)
	SetReadDeadline(t time.Time) error
	Close() error
}

type AuthHandler func(server string, msg *comet.AuthMsg) (ReaderWriterCloser, error)

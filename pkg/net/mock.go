package net

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/txchat/im/api/protocol"
)

type mock struct {
	log     zerolog.Logger
	msgChan chan *protocol.Proto
	hbChan  chan *protocol.Proto
	closer  chan struct{}
}

func NewMockAuth(server string, authMsg *protocol.AuthBody) (ReaderWriterCloser, error) {
	w := io.MultiWriter(os.Stdout)
	log := zerolog.New(w).With().Str("server", server).Str("appId", authMsg.GetAppId()).Str("token", authMsg.GetToken()).Logger()
	c := &mock{
		log:     log,
		msgChan: make(chan *protocol.Proto, 10),
		hbChan:  make(chan *protocol.Proto, 10),
		closer:  make(chan struct{}),
	}
	go c.serve()
	return c, nil
}

func (c *mock) WriteProto(proto *protocol.Proto) (err error) {
	c.log.Info().Interface("proto", proto).Msg("WriteProto")
	switch proto.GetOp() {
	case int32(protocol.Op_Heartbeat):
		var p protocol.Proto
		p.Ver = proto.GetVer()
		p.Op = proto.GetOp()
		p.Seq = proto.GetSeq()
		p.Ack = proto.GetAck()
		p.Body = proto.GetBody()
		select {
		case c.hbChan <- &p:
		default:
		}
	}
	return nil
}

func (c *mock) ReadProto(proto *protocol.Proto) (err error) {
	select {
	case p := <-c.msgChan:
		proto.Ver = p.GetVer()
		proto.Op = p.GetOp()
		proto.Seq = p.GetSeq()
		proto.Ack = p.GetAck()
		proto.Body = p.GetBody()
		c.log.Info().Interface("proto", proto).Msg("ReadProto")
	case <-c.closer:
		return io.EOF
	}
	return nil
}

func (c *mock) SetReadDeadline(t time.Time) error {
	c.log.Info().Interface("time", t).Msg("SetReadDeadline")
	return nil
}

func (c *mock) Close() error {
	close(c.closer)
	c.log.Info().Msg("Close")
	return nil
}

func (c *mock) serve() {
	seq := int32(0)
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			seq++
			p := &protocol.Proto{
				Ver:  0,
				Op:   int32(protocol.Op_Message),
				Seq:  seq,
				Ack:  0,
				Body: nil,
			}
			c.msgChan <- p
		case p := <-c.hbChan:
			seq++
			p.Op = int32(protocol.Op_HeartbeatReply)
			p.Ack = p.Seq
			p.Seq = seq
			c.msgChan <- p
		case <-c.closer:
			return
		}
	}
}

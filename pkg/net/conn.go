package net

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/txchat/im/api/protocol"
)

type ReaderWriterCloser interface {
	WriteProto(proto *protocol.Proto) (err error)
	ReadProto(proto *protocol.Proto) (err error)
	SetReadDeadline(t time.Time) error
	Close() error
}

type AuthHandler func(server string, msg *protocol.AuthBody) (ReaderWriterCloser, error)

type IMConn struct {
	heartTimer *time.Timer
	heartbeat  time.Duration
	connId     string
	seq        int32
	// 发送数据缓存通道
	wb      chan *protocol.Proto
	rb      chan *protocol.Proto
	isClose int32
	closer  chan bool
	conn    ReaderWriterCloser
	err     error
}

// Signal send wb to the channel, protocol ready.
func (c *IMConn) Signal() {
	c.wb <- protocol.ProtoReady
}

// Push server push message.
func (c *IMConn) Push(p *protocol.Proto) (seq int32, err error) {
	p.Seq = c.incSeq()
	seq = p.Seq
	select {
	case c.wb <- p:
	default:
	}
	return
}

// RePush server resend message.
func (c *IMConn) RePush(p *protocol.Proto) (seq int32, err error) {
	seq = p.Seq
	select {
	case c.wb <- p:
	default:
	}
	return
}

func (c *IMConn) Read() *protocol.Proto {
	return <-c.rb
}

// Close the channel.
func (c *IMConn) Close() {
	if atomic.CompareAndSwapInt32(&c.isClose, 0, 1) {
		c.heartTimer.Stop()
		c.err = c.conn.Close()
		close(c.closer)
	}
}

func (c *IMConn) incSeq() int32 {
	return atomic.AddInt32(&c.seq, 1)
}

func (c *IMConn) GetConnId() string {
	return c.connId
}

// dispatch 处理发送操作
func (c *IMConn) dispatch() {
	hbProto := new(protocol.Proto)
	for {
		select {
		case <-c.heartTimer.C:
			// heartbeat
			hbProto.Op = int32(protocol.Op_Heartbeat)
			hbProto.Seq = c.incSeq()
			hbProto.Body = nil

			err := c.conn.WriteProto(hbProto)
			if err != nil {
				c.err = fmt.Errorf("heartbeat WriteProto: %v", err)
				c.Close()
			}
		case <-c.closer:
			return
		case p := <-c.wb:
			switch p {
			case protocol.ProtoFinish:
			case protocol.ProtoReady:
			default:
				err := c.conn.WriteProto(p)
				if err != nil {
					c.err = fmt.Errorf("WriteProto: %v", err)
					c.Close()
				}
			}
		}
	}
}

func (c *IMConn) serve() {
	go c.dispatch()

	// reader
	go func() {
		p := new(protocol.Proto)
		for {
			if err := c.conn.ReadProto(p); err != nil {
				c.err = fmt.Errorf("ReadProto: %v", err)
				c.Close()
				return
			}

			if p.Op == int32(protocol.Op_AuthReply) {
			} else if p.Op == int32(protocol.Op_HeartbeatReply) {
				if err := c.conn.SetReadDeadline(time.Now().Add(c.heartbeat + 60*time.Second)); err != nil {
					c.err = fmt.Errorf("SetReadDeadline: %v", err)
					c.Close()
					return
				}
			} else {
				// 处理从服务端接收的数据
				c.rb <- p
			}
		}
	}()
}

func DialIM(server string, authMsg *protocol.AuthBody, hb time.Duration, auth AuthHandler) (*IMConn, error) {
	// 握手
	conn, err := auth(server, authMsg)
	if err != nil {
		return nil, err
	}
	cli := &IMConn{
		connId:     uuid.New().String(),
		heartTimer: time.NewTimer(hb),
		heartbeat:  hb,
		wb:         make(chan *protocol.Proto, 100),
		rb:         make(chan *protocol.Proto, 100),
		closer:     make(chan bool, 1),
		conn:       conn,
	}
	return cli, nil
}

func DialIMAndServe(server string, authMsg *protocol.AuthBody, hb time.Duration, auth AuthHandler) (*IMConn, error) {
	c, err := DialIM(server, authMsg, hb, auth)
	if err != nil {
		return nil, err
	}
	c.serve()
	return c, nil
}

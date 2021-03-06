package net

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	comet "github.com/txchat/im/api/comet/grpc"
)

type IMConn struct {
	heartbeat time.Duration
	connId    string
	seq       int32
	// 发送数据缓存通道
	wb      chan *comet.Proto
	rb      chan *comet.Proto
	isClose int32
	closer  chan bool
	conn    ReaderWriterCloser
	err     error
}

// Signal send wb to the channel, protocol ready.
func (c *IMConn) Signal() {
	c.wb <- comet.ProtoReady
}

// Push server push message.
func (c *IMConn) Push(p *comet.Proto) (seq int32, err error) {
	p.Seq = c.incSeq()
	seq = p.Seq
	select {
	case c.wb <- p:
	default:
	}
	return
}

// RePush server resend message.
func (c *IMConn) RePush(p *comet.Proto) (seq int32, err error) {
	seq = p.Seq
	select {
	case c.wb <- p:
	default:
	}
	return
}

func (c *IMConn) Read() *comet.Proto {
	return <-c.rb
}

// Close the channel.
func (c *IMConn) Close() {
	if atomic.CompareAndSwapInt32(&c.isClose, 0, 1) {
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
	for {
		select {
		case <-c.closer:
			return
		case p := <-c.wb:
			switch p {
			case comet.ProtoFinish:
			case comet.ProtoReady:
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
		p := new(comet.Proto)
		for {
			if err := c.conn.ReadProto(p); err != nil {
				c.err = fmt.Errorf("ReadProto: %v", err)
				c.Close()
				return
			}

			if p.Op == int32(comet.Op_AuthReply) {
			} else if p.Op == int32(comet.Op_HeartbeatReply) {
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

func DialIM(server string, authMsg *comet.AuthMsg, hb time.Duration, auth AuthHandler) (*IMConn, error) {
	// 握手
	conn, err := auth(server, authMsg)
	if err != nil {
		return nil, err
	}
	cli := &IMConn{
		connId:    uuid.New().String(),
		heartbeat: hb,
		wb:        make(chan *comet.Proto, 100),
		rb:        make(chan *comet.Proto, 100),
		closer:    make(chan bool, 1),
		conn:      conn,
	}

	// 定期发送心跳
	go func() {
		hbProto := new(comet.Proto)
		for {
			// heartbeat
			hbProto.Op = int32(comet.Op_Heartbeat)
			hbProto.Seq = cli.incSeq()
			hbProto.Body = nil
			if _, err := cli.Push(hbProto); err != nil {
				cli.err = err
				return
			}
			//TODO 可以使用二叉堆计时器来优化
			time.Sleep(cli.heartbeat)
			select {
			case <-cli.closer:
				return
			default:
			}
		}
	}()
	return cli, nil
}

func DialIMAndServe(server string, authMsg *comet.AuthMsg, hb time.Duration, auth AuthHandler) (*IMConn, error) {
	c, err := DialIM(server, authMsg, hb, auth)
	if err != nil {
		return nil, err
	}
	c.serve()
	return c, nil
}

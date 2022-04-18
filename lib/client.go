package lib

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	comet "github.com/txchat/im/api/comet/grpc"
)

type Client struct {
	connId  string
	heart   time.Duration
	seq     int32
	signal  chan *comet.Proto
	isClose int32
	closer  chan bool
	biz     Biz
	rw      ReaderWriter

	err error
}

// Ready check the channel ready or close?
func (c *Client) SetBiz(biz Biz) {
	c.biz = biz
}

// Signal send signal to the channel, protocol ready.
func (c *Client) Signal() {
	c.signal <- comet.ProtoReady
}

// Push server push message.
func (c *Client) Push(p *comet.Proto) (err error) {
	select {
	case c.signal <- p:
	default:
	}
	return
}

// Close close the channel.
func (c *Client) Close() {
	if atomic.CompareAndSwapInt32(&c.isClose, 0, 1) {
		close(c.closer)
	}
}

func (c *Client) incSeq() int32 {
	return atomic.AddInt32(&c.seq, 1)
}

func (c *Client) IncSeq() int32 {
	return atomic.AddInt32(&c.seq, 1)
}

func (c *Client) dispatch() {
	for {
		select {
		case <-c.closer:
			return
		case p := <-c.signal:
			switch p {
			case comet.ProtoFinish:
			case comet.ProtoReady:
			default:
				err := c.rw.WriteProto(p)
				if err != nil {
					c.err = fmt.Errorf("WriteProto: %v", err)
					c.Close()
				}
			}
		}
	}
}

func (c *Client) Serve() {
	if c.biz == nil {
		panic("not init biz")
	}
	go c.dispatch()

	// reader
	go func() {
		p := new(comet.Proto)
		for {
			if err := c.rw.ReadProto(p); err != nil {
				c.err = fmt.Errorf("ReadProto: %v", err)
				c.Close()
				return
			}

			//xxx, _ := json.Marshal(p)
			//log.Debug("server resp", "proto", string(xxx))

			if p.Op == int32(comet.Op_AuthReply) {
			} else if p.Op == int32(comet.Op_HeartbeatReply) {
				if err := c.rw.SetReadDeadline(time.Now().Add(c.heart + 60*time.Second)); err != nil {
					c.err = fmt.Errorf("SetReadDeadline: %v", err)
					c.Close()
					return
				}
			} else {
				// biz msg resp
				err := c.biz.Receive(c, p)
				if err != nil {
					c.err = fmt.Errorf("receive callback: %v", err)
					c.Close()
					return
				}
			}
		}
	}()
}

func (c *Client) GetConnId() string {
	return c.connId
}

func NewClient(appId, token, server string, ext []byte, heart time.Duration, auth AuthHandler) (*Client, error) {
	rw, err := auth(appId, token, server, ext)
	if err != nil {
		return nil, err
	}
	connId := uuid.New().String()
	cli := &Client{
		connId:  connId,
		heart:   heart,
		seq:     0,
		signal:  make(chan *comet.Proto, 100),
		isClose: 0,
		closer:  make(chan bool, 1),
		biz:     nil,
		rw:      rw,
	}

	// writer heartbeat
	go func() {
		hbProto := new(comet.Proto)
		for {
			// heartbeat
			hbProto.Op = int32(comet.Op_Heartbeat)
			hbProto.Seq = cli.incSeq()
			hbProto.Body = nil
			if err := cli.Push(hbProto); err != nil {
				cli.err = err
				return
			}
			time.Sleep(cli.heart)
			select {
			case <-cli.closer:
				return
			default:
			}
		}
	}()
	return cli, nil
}

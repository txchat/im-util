package tcp

import (
	"fmt"
	"net"
	"time"

	"github.com/Terry-Mao/goim/pkg/bufio"
	"github.com/golang/protobuf/proto"
	xnet "github.com/txchat/im-util/pkg/net"
	"github.com/txchat/im/api/protocol"
)

type Client struct {
	conn net.Conn
	wr   *bufio.Writer
	rd   *bufio.Reader
}

func Auth(server string, authMsg *protocol.AuthBody) (xnet.ReaderWriterCloser, error) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return nil, fmt.Errorf("auth step 1 dial: %v", err)
	}
	wr := bufio.NewWriter(conn)
	rd := bufio.NewReader(conn)

	p := new(protocol.Proto)
	p.Ver = 1
	p.Op = int32(protocol.Op_Auth)
	p.Seq = 0
	p.Ack = 0
	p.Body, _ = proto.Marshal(authMsg)

	//auth
	if err = p.WriteTCP(wr); err != nil {
		return nil, fmt.Errorf("auth step 2 write auth frame: %v", err)
	}
	// only hungry flush response
	if err = wr.Flush(); err != nil {
		return nil, fmt.Errorf("auth step 2 write auth frame: %v", err)
	}
	if err = p.ReadTCP(rd); err != nil {
		return nil, fmt.Errorf("auth step 3 read auth reply frame: %v", err)
	}
	return &Client{
		conn: conn,
		wr:   wr,
		rd:   rd,
	}, nil
}

func (c *Client) WriteProto(proto *protocol.Proto) (err error) {
	defer func() {
		err = c.wr.Flush()
	}()
	return proto.WriteTCP(c.wr)
}

func (c *Client) ReadProto(proto *protocol.Proto) (err error) {
	return proto.ReadTCP(c.rd)
}

func (c *Client) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *Client) Close() error {
	return c.conn.Close()
}

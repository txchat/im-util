package tcp

import (
	"log"
	"net"
	"time"

	"github.com/Terry-Mao/goim/pkg/bufio"
	"github.com/golang/protobuf/proto"
	"github.com/txchat/im-util/lib"
	comet "github.com/txchat/im/api/comet/grpc"
)

type Client struct {
	conn net.Conn
	wr   *bufio.Writer
	rd   *bufio.Reader
}

func Auth(appId, token, server string) (lib.ReaderWriter, error) {
	// connnect to server
	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Printf("net.Dial(%s) error(%v)", server, err)
		return nil, err
	}
	wr := bufio.NewWriter(conn)
	rd := bufio.NewReader(conn)

	authMsg := &comet.AuthMsg{
		AppId: appId,
		Token: token,
	}
	p := new(comet.Proto)
	p.Ver = 1
	p.Op = int32(comet.Op_Auth)
	p.Seq = 0
	p.Ack = 0
	p.Body, _ = proto.Marshal(authMsg)

	//auth
	if err = p.WriteTCP(wr); err != nil {
		return nil, err
	}
	// only hungry flush response
	if err = wr.Flush(); err != nil {
		return nil, err
	}
	if err = p.ReadTCP(rd); err != nil {
		return nil, err
	}
	return &Client{
		conn: conn,
		wr:   wr,
		rd:   rd,
	}, nil
}

func (c *Client) WriteProto(proto *comet.Proto) (err error) {
	defer func() {
		err = c.wr.Flush()
	}()
	return proto.WriteTCP(c.wr)
}

func (c *Client) ReadProto(proto *comet.Proto) (err error) {
	return proto.ReadTCP(c.rd)
}

func (c *Client) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

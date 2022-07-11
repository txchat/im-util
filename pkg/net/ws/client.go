package ws

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	xnet "github.com/txchat/im-util/pkg/net"
	comet "github.com/txchat/im/api/comet/grpc"
)

type Client struct {
	conn *websocket.Conn
}

func Auth(server string, authMsg *comet.AuthMsg) (xnet.ReaderWriterCloser, error) {
	wsUrl := "ws://" + server + "/sub"
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("auth step 1 dial: %v", err)
	}

	p := new(comet.Proto)
	p.Ver = 1
	p.Op = int32(comet.Op_Auth)
	p.Seq = 0
	p.Ack = 0
	p.Body, _ = proto.Marshal(authMsg)

	//auth
	if err = p.WriteWebsocket2(conn); err != nil {
		return nil, fmt.Errorf("auth step 2 write auth frame: %v", err)
	}
	if err = p.ReadWebsocket2(conn); err != nil {
		return nil, fmt.Errorf("auth step 3 read auth reply frame: %v", err)
	}
	return &Client{conn: conn}, nil
}

func (c *Client) WriteProto(proto *comet.Proto) (err error) {
	return proto.WriteWebsocket2(c.conn)
}

func (c *Client) ReadProto(proto *comet.Proto) (err error) {
	return proto.ReadWebsocket2(c.conn)
}

func (c *Client) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *Client) Close() error {
	return c.conn.Close()
}

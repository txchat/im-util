package ws

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/txchat/im-util/lib"
	comet "github.com/txchat/im/api/comet/grpc"
	"time"
)

type Client struct {
	conn *websocket.Conn
}

func Auth(appId, token, server string, ext []byte) (lib.ReaderWriter, error) {
	// connnect to server
	wsUrl := "ws://" + server + "/sub"
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("auth step 1 dial: %v", err)
	}

	authMsg := &comet.AuthMsg{
		AppId: appId,
		Token: token,
		Ext:   ext,
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
		return nil, fmt.Errorf("auth step 3 reac auth reply frame: %v", err)
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

package user

import (
	"github.com/rs/zerolog/log"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/txchat/dtalk/pkg/auth"
	"github.com/txchat/im-util/lib"
	"github.com/txchat/im-util/lib/ws"
	comet "github.com/txchat/im/api/comet/grpc"
	xproto "github.com/txchat/imparse/proto"
)

type User struct {
	address        string
	userName       string
	priKey, pubKey []byte
	revBuf         chan *comet.Proto

	conn     *lib.Client
	isClosed int32
}

func NewUser(address string, priKey, pubKey []byte) *User {
	return &User{
		address:  address,
		priKey:   priKey,
		pubKey:   pubKey,
		revBuf:   make(chan *comet.Proto, 100),
		isClosed: 1,
	}
}

func (u *User) AuthToken() string {
	authenticator := auth.NewDefaultApuAuthenticator()
	return authenticator.Request("", u.pubKey, u.priKey)
}

func (u *User) Close() {
	if atomic.CompareAndSwapInt32(&u.isClosed, 0, 1) {
		u.conn.Close()
		u.conn = nil
	}
	log.Info().Msg("User Close")
}

func (u *User) ConnServer(appId, server string) error {
	if atomic.CompareAndSwapInt32(&u.isClosed, 1, 0) {
		cli, err := lib.NewClient(appId, u.AuthToken(), server, nil, 20*time.Second, ws.Auth)
		if err != nil {
			return err
		}
		cli.SetBiz(u)
		cli.Serve()
		u.conn = cli
	}
	return nil
}

func (u *User) ConnServerWithDevice(appId, server string, devInfo *xproto.Login) error {
	if atomic.CompareAndSwapInt32(&u.isClosed, 1, 0) {
		extData, _ := proto.Marshal(devInfo)

		cli, err := lib.NewClient(appId, u.AuthToken(), server, extData, 20*time.Second, ws.Auth)
		if err != nil {
			return err
		}
		cli.SetBiz(u)
		cli.Serve()
		u.conn = cli
	}
	return nil
}

func (u *User) Receive(c *lib.Client, proto *comet.Proto) error {
	u.revBuf <- proto
	return nil
}

func (u *User) OnReceive() *comet.Proto {
	return <-u.revBuf
}

func (u *User) Send(proto *comet.Proto) error {
	return u.conn.Push(proto)
}

func (u *User) GetId() string {
	return u.address
}

func (u *User) GetConnId() string {
	return u.conn.GetConnId()
}

func (u *User) GenSeq() int32 {
	return u.conn.IncSeq()
}

func (u *User) SetUsername(username string) {
	u.userName = username
}

func (u *User) GetUsername() string {
	if u.userName == "" {
		return u.address
	}
	return u.userName
}

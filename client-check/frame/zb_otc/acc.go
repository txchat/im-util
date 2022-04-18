package dtalk

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/txchat/im-util/client-check/model"
	comet "github.com/txchat/im/api/comet/grpc"
)

type account struct {
}

func (t *account) Check(p *comet.Proto) (err error) {
	var (
		auth comet.AuthMsg
	)
	err = proto.Unmarshal(p.Body, &auth)
	if err != nil {
		return
	}

	if auth.AppId == "" || auth.Token == "" {
		err = model.ErrInvalidAuthReq
		return
	}

	if !t.checkAppId(auth.AppId) {
		fmt.Errorf("invalid appId %s", auth.AppId)
		err = model.ErrInvalidAppId
		return
	}
	if auth.Token == "" {
		err = model.ErrInvalidToken
	}
	return
}

func (t *account) Set(key string, val interface{}) {
}

func (t *account) checkAppId(input string) bool {
	switch input {
	case "zb_otc":
		return true
	default:
		return false
	}
}

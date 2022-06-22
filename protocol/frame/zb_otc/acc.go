package dtalk

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/txchat/dtalk/service/auth/model"
	"github.com/txchat/im-util/protocol/frame"
	comet "github.com/txchat/im/api/comet/grpc"
)

type authChecker struct {
}

func (t *authChecker) Check(p *comet.Proto) (err error) {
	var (
		authFrame comet.AuthMsg
	)
	err = proto.Unmarshal(p.Body, &authFrame)
	if err != nil {
		return
	}

	if authFrame.AppId == "" || authFrame.Token == "" {
		err = frame.ErrInvalidAuthReq
		return
	}

	if !t.checkAppId(authFrame.AppId) {
		err = fmt.Errorf("%v: %s", frame.ErrInvalidAppId, authFrame.AppId)
		return
	}
	if authFrame.Token == "" {
		err = model.ErrInvalidToken
	}
	return
}

func (t *authChecker) Set(key string, val interface{}) {
}

func (t *authChecker) checkAppId(input string) bool {
	switch input {
	case "zb_otc":
		return true
	}
	return false
}

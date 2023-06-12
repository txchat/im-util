package dtalk

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/txchat/dtalk/pkg/auth"
	"github.com/txchat/im-util/protocol/frame"
	"github.com/txchat/im/api/protocol"
)

type authChecker struct {
}

func (t *authChecker) Check(p *protocol.Proto) (err error) {
	var (
		authFrame protocol.AuthBody
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
		err = auth.ErrSignatureInvalid(err)
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

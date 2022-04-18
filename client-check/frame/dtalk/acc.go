package dtalk

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/txchat/dtalk/pkg/address"
	bizApi "github.com/txchat/dtalk/pkg/api"
	xerror "github.com/txchat/dtalk/pkg/error"
	"github.com/txchat/im-util/client-check/model"
	comet "github.com/txchat/im/api/comet/grpc"
)

type account struct {
	isCkTimeOut bool
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
	pubKey, err := bizApi.VerifyAddress(auth.Token)
	if err != nil {
		if model.ParseErr(err) != xerror.SignatureExpired || t.isCkTimeOut {
			return
		}
		err = nil
	}
	addr := address.PublicKeyToAddress(address.NormalVer, pubKey)
	if addr == "" {
		return errors.New("pubkey can not to address")
	}
	return
}

func (t *account) Set(key string, val interface{}) {
	switch key {
	case "isCkTimeOut":
		switch val.(type) {
		case bool:
			t.isCkTimeOut = val.(bool)
		default:
			fmt.Printf("%T\n", val)
		}
	default:
		fmt.Println(key)
	}
}

func (t *account) checkAppId(input string) bool {
	switch input {
	case "dtalk":
		return true
	default:
		return false
	}
}

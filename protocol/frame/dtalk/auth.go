package dtalk

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/txchat/dtalk/pkg/address"
	"github.com/txchat/dtalk/pkg/auth"
	xcrypt "github.com/txchat/dtalk/pkg/crypt"

	//secp256k1_ethereum "github.com/txchat/dtalk/pkg/crypt/secp256k1-ethereum"
	secp256k1_haltingstate "github.com/txchat/dtalk/pkg/crypt/secp256k1-haltingstate"
	"github.com/txchat/im-util/protocol/frame"
	"github.com/txchat/im/api/protocol"
)

var driver xcrypt.Encrypt

func init() {
	var err error
	driver, err = xcrypt.Load(secp256k1_haltingstate.Name)
	if err != nil {
		panic(err)
	}
}

type authChecker struct {
	isCkTimeOut bool
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

	if !t.checkAppID(authFrame.AppId) {
		err = fmt.Errorf("%v: %s", frame.ErrInvalidAppId, authFrame.AppId)
		return
	}

	apiRequest, err := auth.NewAPIRequestFromToken(authFrame.Token)
	if err != nil {
		return err
	}
	signatory, err := auth.NewSignatoryFromMetadata(driver, apiRequest.GetMetadata())
	if err != nil {
		return err
	}
	if isMatch, err := signatory.Match(apiRequest.GetSignature(), apiRequest.GetPublicKey()); !isMatch {
		return auth.ErrSignatureInvalid(err)
	}
	if t.isCkTimeOut && signatory.IsExpire() {
		return auth.ErrSignatureExpired
	}
	uid := address.PublicKeyToAddress(address.NormalVer, apiRequest.GetPublicKey())
	if uid == "" {
		return auth.ErrUIDInvalid
	}
	return
}

func (t *authChecker) Set(key string, val interface{}) {
	switch key {
	case "isCkTimeOut":
		switch val := val.(type) {
		case bool:
			t.isCkTimeOut = val
		default:
			fmt.Printf("%T\n", val)
		}
	default:
		fmt.Println(key)
	}
}

func (t *authChecker) checkAppID(input string) bool {
	switch input {
	case "dtalk":
		return true
	}
	return false
}

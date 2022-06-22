package dtalk

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/txchat/dtalk/pkg/address"
	"github.com/txchat/dtalk/pkg/auth"
	xcrypt "github.com/txchat/dtalk/pkg/crypt"
	secp256k1_ethereum "github.com/txchat/dtalk/pkg/crypt/secp256k1-ethereum"
	"github.com/txchat/im-util/protocol/frame"
	comet "github.com/txchat/im/api/comet/grpc"
)

var driver xcrypt.Encrypt

func init() {
	var err error
	driver, err = xcrypt.Load(secp256k1_ethereum.Name)
	if err != nil {
		panic(err)
	}
}

type authChecker struct {
	isCkTimeOut bool
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

	apiRequest, err := auth.NewApiRequestFromToken(authFrame.Token)
	if err != nil {
		return err
	}
	signatory, err := auth.NewSignatoryFromMetadata(driver, apiRequest.GetMetadata())
	if err != nil {
		return err
	}
	if isMatch, err := signatory.Match(apiRequest.GetSignature(), apiRequest.GetPublicKey()); !isMatch {
		return auth.ERR_SIGNATUREINVALID(err)
	}
	if t.isCkTimeOut && signatory.IsExpire() {
		return auth.ERR_SIGNATUREEXPIRED
	}
	uid := address.PublicKeyToAddress(address.NormalVer, apiRequest.GetPublicKey())
	if uid == "" {
		return auth.ERR_UIDINVALID
	}
	return
}

func (t *authChecker) Set(key string, val interface{}) {
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

func (t *authChecker) checkAppId(input string) bool {
	switch input {
	case "dtalk":
		return true
	}
	return false
}

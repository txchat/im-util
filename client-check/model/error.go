package model

import (
	"errors"
	xerror "github.com/txchat/dtalk/pkg/error"
)

var (
	// ErrProtoPackLen proto packet len error
	ErrProtoPackLen = errors.New("default server codec pack length error")
	// ErrProtoHeaderLen proto header len error
	ErrProtoHeaderLen = errors.New("default server codec header length error")

	ErrInvalidAuthReq = errors.New("ErrInvalidAuthReq")
	ErrInvalidAppId   = errors.New("ErrInvalidAppId")
	ErrInvalidToken   = errors.New("ErrInvalidToken")
	ErrInvalidPlugin  = errors.New("ErrInvalidAppId:plugin unregister")

	ErrPrivateKeyErr = errors.New("private key is not match")
	ErrPublicKeyErr  = errors.New("public key is not match")
	ErrAddressErr    = errors.New("address key is not match")
)

func ParseErr(err error) int {
	switch e := err.(type) {
	case *xerror.Error:
		return e.Code()
	default:
		return 0
	}
}

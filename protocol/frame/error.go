package frame

import "errors"

var (
	// ErrProtoPackLen proto packet len error
	ErrProtoPackLen = errors.New("default server codec pack length error")
	// ErrProtoHeaderLen proto header len error
	ErrProtoHeaderLen = errors.New("default server codec header length error")

	ErrInvalidAuthReq = errors.New("ErrInvalidAuthReq")
	ErrInvalidAppId   = errors.New("ErrInvalidAppId")
	ErrInvalidToken   = errors.New("ErrInvalidToken")

	ErrPrivateKeyErr = errors.New("private key is not match")
	ErrPublicKeyErr  = errors.New("public key is not match")
	ErrAddressErr    = errors.New("address key is not match")

	ErrInvalidPlugin = errors.New("ErrInvalidAppId:plugin unregister")
)

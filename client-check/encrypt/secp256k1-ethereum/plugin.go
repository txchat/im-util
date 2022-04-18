package secp256K1

import (
	"github.com/txchat/im-util/client-check/encrypt"
)

const Name = "secp256k1-ethereum"

func init() {
	encrypt.Register(Name, New())
}

func New() encrypt.Encrypt {
	return &ethereum{}
}

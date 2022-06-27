package dtalk

import (
	"github.com/txchat/im-util/protocol/frame"
)

const Name = "zb_otc"

func init() {
	frame.RegisterAuthChecker(Name, NewAuthChecker())
}

func NewAuthChecker() frame.AuthChecker {
	return &authChecker{}
}

package dtalk

import (
	"github.com/txchat/im-util/client-check/frame"
)

const Name = "zb_otc"

func init() {
	frame.Register(Name, New())
}

func New() frame.Account {
	return &account{}
}

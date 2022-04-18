package dtalk

import (
	"github.com/txchat/im-util/client-check/frame"
)

const Name = "dtalk"

func init() {
	frame.Register(Name, New())
}

func New() frame.Account {
	return &account{}
}

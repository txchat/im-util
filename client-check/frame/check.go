package frame

import (
	"github.com/txchat/im-util/client-check/model"
	comet "github.com/txchat/im/api/comet/grpc"
)

var accountFactory = make(map[string]Account)

func Register(name string, exec Account) {
	accountFactory[name] = exec
}

func Load(name string) (Account, error) {
	exec, ok := accountFactory[name]
	if !ok {
		return nil, model.ErrInvalidPlugin
	}
	return exec, nil
}

type Account interface {
	Set(key string, val interface{})
	Check(p *comet.Proto) error
}

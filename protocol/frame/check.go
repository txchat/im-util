package frame

import (
	"github.com/txchat/im/api/protocol"
)

var authCheckerFactory = make(map[string]AuthChecker)

func RegisterAuthChecker(name string, exec AuthChecker) {
	authCheckerFactory[name] = exec
}

func LoadAuthChecker(name string) (AuthChecker, error) {
	exec, ok := authCheckerFactory[name]
	if !ok {
		return nil, ErrInvalidPlugin
	}
	return exec, nil
}

type AuthChecker interface {
	Set(key string, val interface{})
	Check(p *protocol.Proto) error
}

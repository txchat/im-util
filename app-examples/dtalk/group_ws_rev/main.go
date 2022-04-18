package main

// ./client {appId} {token} {ws-server-addr}
// ./client echo f3dc8ccd localhost:3102

import (
	"flag"
	"github.com/inconshreveable/log15"
	"os"
	"runtime"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/txchat/im-util/lib"
	"github.com/txchat/im-util/lib/ws"
	comet "github.com/txchat/im/api/comet/grpc"
)

var log = log15.New()
var Conf *Config

//os.Args[1]=appId;os.Args[1]=token;os.Args[1]=server
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	_, err := toml.DecodeFile(os.Args[1], &Conf)
	if err != nil {
		panic(err)
	}
	log.Debug("get users", "num", len(Conf.Users))

	for i, user := range Conf.Users {
		log.Debug("create user", "num", i, "token", user.Token, "Uid", user.Uid, "Groups", user.Groups)
		go client(Conf.AppId, user.Token, Conf.Comet, Conf.Logic, user.Uid, user.Groups)
	}
	var exit chan bool
	<-exit
}

func client(appId, token, server, logicAddr, uid string, groups []string) {
	cli, err := lib.NewClient(appId, token, server, 5*time.Second, ws.Auth)
	if err != nil {
		panic(err)
	}
	cli.SetBiz(new(biz))
	cli.Serve()
	g := lib.NewGroup(appId, uid, logicAddr)
	err = g.JoinIn(groups)
	if err != nil {
		log.Error("JoinIn", "err", err)
		panic(err)
	}
}

type biz struct {
}

func (b *biz) Receive(c *lib.Client, p *comet.Proto) error {
	log.Debug("get msg", "uid", c.GetUid(), "proto", p)
	return nil
}

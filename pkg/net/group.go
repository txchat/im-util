package net

import (
	"context"

	"github.com/txchat/im/app/logic/logicclient"
	"github.com/zeromicro/go-zero/zrpc"
)

type Group struct {
	appId       string
	uid         string
	logicClient logicclient.Logic
}

func NewGroup(appId, uid, addr string) *Group {
	return &Group{
		appId: appId,
		uid:   uid,
		logicClient: logicclient.NewLogic(zrpc.MustNewClient(zrpc.RpcClientConf{
			Endpoints: []string{addr},
		}, zrpc.WithNonBlock())),
	}
}

func (g *Group) JoinInByKey(key string, gids []string) error {
	msg := &logicclient.JoinGroupByKeyReq{
		AppId: g.appId,
		Key:   []string{key},
		Gid:   gids,
	}
	_, err := g.logicClient.JoinGroupByKey(context.Background(), msg)
	return err
}

func (g *Group) JoinIn(gids []string) error {
	msg := &logicclient.JoinGroupByUIDReq{
		AppId: g.appId,
		Uid:   []string{g.uid},
		Gid:   gids,
	}
	_, err := g.logicClient.JoinGroupByUID(context.Background(), msg)
	return err
}

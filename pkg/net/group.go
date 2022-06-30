package net

import (
	"context"
	logic "github.com/txchat/im/api/logic/grpc"
	"github.com/txchat/im/common"
	"time"
)

type Group struct {
	appId       string
	uid         string
	logicClient logic.LogicClient
}

func NewGroup(appId, uid, addr string) *Group {
	return &Group{
		appId:       appId,
		uid:         uid,
		logicClient: newLogicClient(addr),
	}
}

func (g *Group) JoinInByKey(key string, gids []string) error {
	msg := &logic.GroupsKey{
		AppId: g.appId,
		Keys:  []string{key},
		Gid:   gids,
	}
	_, err := g.logicClient.JoinGroupsByKeys(context.Background(), msg)
	return err
}

func (g *Group) JoinIn(gids []string) error {
	msg := &logic.GroupsMid{
		AppId: g.appId,
		Mids:  []string{g.uid},
		Gid:   gids,
	}
	_, err := g.logicClient.JoinGroupsByMids(context.Background(), msg)
	return err
}

func newLogicClient(addr string) logic.LogicClient {
	conn, err := common.NewGRPCConn(addr, time.Second)
	if err != nil {
		panic(err)
	}
	return logic.NewLogicClient(conn)
}

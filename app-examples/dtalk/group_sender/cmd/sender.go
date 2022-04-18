package cmd

import (
	"context"
	"github.com/txchat/im-util/app-examples/dtalk/group_sender/config"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/txchat/im-util/app-examples/dtalk/group_sender/internal"
	comet "github.com/txchat/im/api/comet/grpc"
	logic "github.com/txchat/im/api/logic/grpc"
	"github.com/txchat/im/common"
	bizProto "github.com/txchat/imparse/proto"
)

func init() {
	sendCmd.Flags().StringVarP(&config.Conf.LogicAddr, "address", "a", "", "logic server address")
	sendCmd.Flags().StringVarP(&config.Conf.Group, "group", "g", "", "send group")
	rootCmd.AddCommand(sendCmd)
}

func newLogicClient(addr string) logic.LogicClient {
	conn, err := common.NewGRPCConn(addr, time.Second)
	if err != nil {
		panic(err)
	}
	return logic.NewLogicClient(conn)
}

var sendCmd = &cobra.Command{
	Use:     "send",
	Short:   "send groups msg",
	Long:    "",
	Example: "send -a 172.16.101.127:3013 -g 1",
	Run:     SendGroupsMsg,
}

func SendGroupsMsg(cmd *cobra.Command, args []string) {
	logicClient := newLogicClient(config.Conf.LogicAddr)
	biz := internal.TextMsg(bizProto.ToGroup, "1", config.Conf.Group, "hello")
	bizData, err := proto.Marshal(biz)
	if err != nil {
		log.Err(err).Msg("bizData Marshal error")
		return
	}
	p := comet.Proto{
		Ver:  0,
		Op:   int32(comet.Op_ReceiveMsg),
		Seq:  0,
		Ack:  0,
		Body: bizData,
	}
	data, err := proto.Marshal(&p)
	if err != nil {
		log.Err(err).Msg("Proto Marshal error")
		return
	}
	_, err = logicClient.PushGroup(context.Background(), &logic.GroupMsg{
		AppId: "dtalk",
		Group: config.Conf.Group,
		Type:  0,
		Op:    0,
		Msg:   data,
	})
	if err != nil {
		log.Err(err).Interface("config", config.Conf).Msg("PushGroup error")
		return
	}
}

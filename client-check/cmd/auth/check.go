package auth

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"
	"github.com/txchat/im-util/client-check/frame"
	comet "github.com/txchat/im/api/comet/grpc"
)

var (
	appType string
	data    string
)

var CheckCmd = &cobra.Command{
	Use:     "check",
	Short:   "check auth header",
	Long:    "",
	Example: "check -T=dtalk -d fffff ",
	Run:     doCheckCmd,
}

func init() {
	CheckCmd.Flags().StringVarP(&data, "data", "d", "", "the frame data encoded by base64")
	CheckCmd.Flags().StringVarP(&appType, "type", "T", "dtalk", "check app type -T=dtalk")
}

func doCheckCmd(cmd *cobra.Command, args []string) {
	var p comet.Proto
	authF := &comet.AuthMsg{
		AppId: appType,
		Token: data,
		Ext:   nil,
	}
	d, err := proto.Marshal(authF)
	if err != nil {
		fmt.Printf("AuthMsg proto Marshal failed err:%v\n", err)
		return
	}
	p.Body = d

	acc, err := frame.Load(appType)
	if err != nil {
		fmt.Printf("addressCheck failed err:%v\n", err)
		return
	}
	err = acc.Check(&p)
	if err != nil {
		fmt.Printf("addressCheck failed err:%v\n", err)
		return
	}
	fmt.Printf("auth frame check success!")
	return
}

package auth

import (
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/txchat/im-util/client-check/frame"
	_ "github.com/txchat/im-util/client-check/frame/dtalk"
	_ "github.com/txchat/im-util/client-check/frame/zb_otc"
)

var (
	dataBase64  string
	isCkTimeOut bool
)

var CheckCmd = &cobra.Command{
	Use:     "check",
	Short:   "check im auth frame",
	Long:    "check im auth frame",
	Example: " imall auth -T=dtalk -t=false -d <base64(data)>",
	Run:     doCheckCmd,
}

func init() {
	CheckCmd.Flags().StringVarP(&dataBase64, "data", "d", "", "the frame data encoded by base64")
	CheckCmd.Flags().StringVarP(&appType, "type", "T", "dtalk", "check app type -T=[dtalk]")
	CheckCmd.Flags().BoolVarP(&isCkTimeOut, "timeout", "t", true, "check timeout enable -t=[true]")
}

func doCheckCmd(cmd *cobra.Command, args []string) {
	data, err := base64.StdEncoding.DecodeString(dataBase64)
	if err != nil {
		fmt.Printf("input illegal:%v\n", dataBase64)
		return
	}
	p, err := frame.ToProto(data)
	if err != nil {
		fmt.Printf("frameCheck failed err:%v\n", err)
		return
	}
	//
	acc, err := frame.Load(appType)
	if err != nil {
		fmt.Printf("addressCheck failed err:%v\n proto:%v\n", err, p)
		return
	}
	acc.Set("isCkTimeOut", isCkTimeOut)
	err = acc.Check(&p)
	if err != nil {
		fmt.Printf("addressCheck failed err:%v\n proto:%v\n", err, p)
		return
	}
	fmt.Printf("auth frame check success!")
	return
}

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
	dataBase64 string
)

var CheckCmd = &cobra.Command{
	Use:     "check",
	Short:   "check im frame",
	Long:    "check im frame",
	Example: " imall frame check -d <base64(data)>",
	Run:     doCheckCmd,
}

func init() {
	CheckCmd.Flags().StringVarP(&dataBase64, "data", "d", "", "the frame data encoded by base64")
}

func doCheckCmd(cmd *cobra.Command, args []string) {
	data, err := base64.StdEncoding.DecodeString(dataBase64)
	if err != nil {
		fmt.Printf("input illegal:%v\n", dataBase64)
		return
	}
	_, err = frame.ToProto(data)
	if err != nil {
		fmt.Printf("frameCheck failed err:%v\n", err)
		return
	}
	fmt.Printf("frame check success!")
	return
}

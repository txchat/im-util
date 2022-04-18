package auth

import (
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"
	"github.com/txchat/im-util/client-check/frame"
	comet "github.com/txchat/im/api/comet/grpc"
)

var GenCmd = &cobra.Command{
	Use:     "gen",
	Short:   "generate auth frame",
	Long:    "auth token param is not encoded. auth frame encoded by base64.This result can used by imall auth command",
	Example: "  imall gen auth -T dtalk -t [token]",
	Run:     doGenCmd,
}

var (
	ver  int32
	seq  int32
	ack  int32
	body []byte
)

var (
	appType string
	token   string
	ext     string
)

func init() {
	GenCmd.Flags().Int32VarP(&ver, "version", "v", 0, "proto version field")
	GenCmd.Flags().Int32VarP(&seq, "seq", "s", 0, "proto seq field")
	GenCmd.Flags().Int32VarP(&ack, "ack", "a", 0, "proto ack field")

	GenCmd.Flags().StringVarP(&appType, "type", "T", "dtalk", "app type -T=[dtalk]")
	GenCmd.Flags().StringVarP(&token, "token", "t", "dtalk", "auth token")
}

func doGenCmd(cmd *cobra.Command, args []string) {
	authF := &comet.AuthMsg{
		AppId: appType,
		Token: token,
		Ext:   nil,
	}
	body, err := proto.Marshal(authF)
	if err != nil {
		fmt.Printf("generate failed err:%v\n", err)
		return
	}
	p := comet.Proto{
		Ver:  ver,
		Op:   int32(comet.Op_Auth),
		Seq:  seq,
		Ack:  ack,
		Body: body,
	}
	data := frame.ToBytes(p)
	fmt.Println(base64.StdEncoding.EncodeToString(data))
}

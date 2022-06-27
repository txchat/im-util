/*
Copyright © 2022 oofpgDLD <oofpgdld@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"github.com/Terry-Mao/goim/pkg/bufio"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"
	comet "github.com/txchat/im/api/comet/grpc"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "生成IM鉴权帧数据",
	Long: `

鉴权数据帧输出格式为Base64([]byte)
`,
	RunE: genRunE,
}

var (
	ver int32
	seq int32
	ack int32

	appType string
	token   string
	ext     string
)

func init() {
	Cmd.AddCommand(genCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	genCmd.Flags().Int32VarP(&ver, "version", "v", 0, "proto version field")
	genCmd.Flags().Int32VarP(&seq, "seq", "s", 0, "proto seq field")
	genCmd.Flags().Int32VarP(&ack, "ack", "a", 0, "proto ack field")

	genCmd.Flags().StringVarP(&appType, "type", "T", "dtalk", "app type -T=[dtalk]")
	genCmd.Flags().StringVarP(&token, "token", "t", "", "auth token")

	genCmd.MarkFlagRequired("token")
}

func genRunE(cmd *cobra.Command, args []string) error {
	var extData []byte
	var err error
	if ext != "" {
		extData, err = hex.DecodeString(ext)
		if err != nil {
			return err
		}
	}

	authFrame := &comet.AuthMsg{
		AppId: appType,
		Token: token,
		Ext:   extData,
	}
	authFrameData, err := proto.Marshal(authFrame)
	if err != nil {
		return err
	}

	p := &comet.Proto{
		Ver:  ver,
		Op:   int32(comet.Op_Auth),
		Seq:  seq,
		Ack:  ack,
		Body: authFrameData,
	}

	buffer := new(bytes.Buffer)

	wr := bufio.NewWriter(buffer)
	err = p.WriteTCP(wr)
	if err != nil {
		return err
	}
	err = wr.Flush()
	if err != nil {
		return err
	}
	buf := buffer.Bytes()
	cmd.Printf("base64 encoding:%s\n", base64.StdEncoding.EncodeToString(buf))
	cmd.Printf("hex encoding:%s\n", hex.EncodeToString(buf))
	return nil
}

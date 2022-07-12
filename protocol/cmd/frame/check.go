// Package frame
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
package frame

import (
	"bytes"
	"encoding/base64"

	"github.com/Terry-Mao/goim/pkg/bufio"
	"github.com/spf13/cobra"
	comet "github.com/txchat/im/api/comet/grpc"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "校验IM数据帧格式是否正确",
	Long: `只进行基本格式的校验，具体的内部字段不进项校验。如果需要具体某种类型的数据帧，请参考：
鉴权帧校验--使用子命令「auth check」
`,
	RunE: checkRunE,
}

func init() {
	Cmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	Cmd.Flags().StringVarP(&dataBase64, "data", "d", "", "the frame data encoded by base64")
	Cmd.MarkFlagRequired("data")
}

var (
	dataBase64 string
)

func checkRunE(cmd *cobra.Command, args []string) error {
	data, err := base64.StdEncoding.DecodeString(dataBase64)
	if err != nil {
		return err
	}

	var p comet.Proto
	err = p.ReadTCP(bufio.NewReader(bytes.NewReader(data)))
	if err != nil {
		return err
	}
	cmd.Println("frame frame legitimate!")
	return nil
}

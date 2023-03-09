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

	"github.com/Terry-Mao/goim/pkg/bufio"
	"github.com/spf13/cobra"
	"github.com/txchat/im-util/protocol/frame"
	"github.com/txchat/im/api/protocol"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "校验鉴权帧",
	Long:  `通过标签-d将鉴权帧的完整数据内容传递给命令工具来校验鉴权是否通过，可以选择是否开启过期时间。`,
	RunE:  checkRunE,
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
	checkCmd.Flags().StringVarP(&dataBase64, "data", "d", "", "the frame data encoded by base64")
	checkCmd.Flags().StringVarP(&appType, "type", "T", "dtalk", "check app type -T=[dtalk]")
	checkCmd.Flags().BoolVarP(&isCkTimeOut, "timeout", "", false, "check timeout enable -t=[false]")

	checkCmd.MarkFlagRequired("data")
}

var (
	dataBase64  string
	isCkTimeOut bool
)

func checkRunE(cmd *cobra.Command, args []string) error {
	data, err := base64.StdEncoding.DecodeString(dataBase64)
	if err != nil {
		return err
	}

	var p protocol.Proto
	err = p.ReadTCP(bufio.NewReader(bytes.NewReader(data)))
	if err != nil {
		return err
	}

	checker, err := frame.LoadAuthChecker(appType)
	if err != nil {
		return err
	}
	checker.Set("isCkTimeOut", isCkTimeOut)
	err = checker.Check(&p)
	if err != nil {
		return err
	}
	cmd.Println("auth frame check success!")
	return nil
}

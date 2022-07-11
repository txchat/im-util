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
package signature

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"

	"github.com/spf13/cobra"
	xcrypt "github.com/txchat/dtalk/pkg/crypt"
)

// sigCmd represents the sig command
var sigCmd = &cobra.Command{
	Use:   "sig",
	Short: "签名一段给定的原文",
	Long: `指定一种签名算法，使用该算法签名sha256后的原文数据

默认的签名算法库「secp256k1_haltingstate」
目前支持：
1. 「secp256k1_haltingstate」
2. 「secp256k1_ethereum」`,
	RunE: sigRunE,
}

func init() {
	Cmd.AddCommand(sigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	sigCmd.Flags().StringVarP(&privateKey, "private", "p", "", "私钥的十六进制字符串编码")
	sigCmd.Flags().IntVarP(&algorithm, "type", "t", 1, "签名算法类型")
	sigCmd.Flags().StringVarP(&msg, "msg", "m", "", "待签名原文")

	sigCmd.MarkFlagRequired("private")
	sigCmd.MarkFlagRequired("msg")
}

func sigRunE(cmd *cobra.Command, args []string) error {
	private, err := hex.DecodeString(privateKey)
	if err != nil {
		return err
	}
	alName := algorithmName[algorithmType(algorithm)]
	if alName == "" {
		return errors.New("algorithm not support")
	}
	driver, err := xcrypt.Load(alName)
	if err != nil {
		return err
	}
	msg256 := sha256.Sum256([]byte(msg))
	data, err := driver.Sign(msg256[:], private)
	if err != nil {
		return err
	}
	cmd.Printf("hex format: %s\n", hex.EncodeToString(data))
	cmd.Printf("base64 format: %s\n", base64.StdEncoding.EncodeToString(data))
	return nil
}

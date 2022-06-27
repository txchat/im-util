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
	"encoding/hex"
	"errors"
	"github.com/spf13/cobra"
	xcrypt "github.com/txchat/dtalk/pkg/crypt"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "A brief description of your command",
	Long: `指定一种签名算法，使用该算法验签

默认的签名算法库「secp256k1_haltingstate」
目前支持：
1. 「secp256k1_haltingstate」
2. 「secp256k1_ethereum」`,
	RunE: verifyRunE,
}

func init() {
	Cmd.AddCommand(verifyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// verifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// verifyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	verifyCmd.Flags().StringVarP(&publicKey, "public", "P", "", "公钥的十六进制字符串编码")
	verifyCmd.Flags().IntVarP(&algorithm, "type", "t", 1, "签名算法类型")
	verifyCmd.Flags().StringVarP(&msg, "msg", "m", "", "待签名原文")
	verifyCmd.Flags().StringVarP(&sig, "sig", "s", "", "签名后数据，十六进制格式")

	verifyCmd.MarkFlagRequired("public")
	verifyCmd.MarkFlagRequired("msg")
	verifyCmd.MarkFlagRequired("sig")
}

func verifyRunE(cmd *cobra.Command, args []string) error {
	public, err := hex.DecodeString(publicKey)
	if err != nil {
		return err
	}
	sigData, err := hex.DecodeString(sig)
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
	b, err := driver.Verify(msg256[:], sigData, public)
	if err != nil {
		return err
	}
	cmd.Printf("verify result is %v\n", b)
	return nil
}

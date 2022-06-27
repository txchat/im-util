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
package dh

import (
	"encoding/hex"
	"strings"

	walletapi "github.com/txchat/chatcipher"

	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:     "encrypt",
	Short:   "通过会话秘钥或者公私钥加密元数据",
	Long:    ``,
	Example: "encrypt -d '' -k ''",
	Run:     encryptRun,
}

func init() {
	Cmd.AddCommand(encryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	encryptCmd.Flags().StringVarP(&metadata, "metadata", "d", "", "待加密元数据，hex编码")
	encryptCmd.Flags().StringVarP(&sessionKey, "key", "k", "", "加密会话秘钥，由private key和public key生成，hex编码")
	encryptCmd.Flags().StringVarP(&privateKey, "private", "p", "", "加密私钥，hex编码")
	encryptCmd.Flags().StringVarP(&publicKey, "public", "P", "", "加密公钥，hex编码")

	encryptCmd.MarkFlagsRequiredTogether("private", "public")
	encryptCmd.MarkFlagRequired("metadata")
}

func encryptRun(cmd *cobra.Command, args []string) {
	privateKey = strings.Replace(privateKey, "0x", "", 1)
	publicKey = strings.Replace(publicKey, "0x", "", 1)
	sessionKey = strings.Replace(sessionKey, "0x", "", 1)
	data, err := hex.DecodeString(strings.Replace(metadata, "0x", "", 1))
	if err != nil {
		cmd.PrintErr("hex.DecodeString metadata failed:%v\n", err)
		return
	}

	//优先秘钥对加密
	if privateKey != "" && publicKey != "" {
		encryptedData, err := walletapi.EncryptWithDHKeyPair(privateKey, publicKey, data)
		if err != nil {
			cmd.PrintErr("walletapi.EncryptWithDHKeyPair failed err:%v\n", err)
			return
		}
		cmd.Printf("encrypt by dh key pair success!:%s\n", hex.EncodeToString(encryptedData))
	} else {
		encryptedData, err := walletapi.EncryptSymmetric(sessionKey, data)
		if err != nil {
			cmd.PrintErr("walletapi.EncryptSymmetric failed err:%v\n", err)
			return
		}
		cmd.Printf("encrypt by dh session key success!:%s\n", hex.EncodeToString(encryptedData))
	}
	return
}

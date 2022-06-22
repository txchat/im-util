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
package token

import (
	"encoding/hex"
	"github.com/spf13/cobra"
	"github.com/txchat/dtalk/pkg/auth"
	"github.com/txchat/im-util/protocol/wallet"
	"time"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "生成im鉴权token",
	Long: `通过助记词或者秘钥对生成指定App下的鉴权Token。
签名和验签的过程参考：https://github.com/txchat/dtalk/tree/main/pkg/auth
`,
	Example: `gen -m '游 即 暗 体 柬 京 非 李 限 稻 跳 务 桥 凶 溶' -a 'dtalk'`,
	RunE:    genRunE,
}

var (
	privateKey string
	publicKey  string
	mnemonic   string
	timestamp  int64
	appKey     string
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
	genCmd.Flags().StringVarP(&privateKey, "private", "p", "", "私钥的十六进制字符串编码")
	genCmd.Flags().StringVarP(&publicKey, "public", "P", "", "公钥的十六进制字符串编码")
	genCmd.Flags().Int64VarP(&timestamp, "timestamp", "t", time.Now().UnixMilli(), "签名时间戳（默认，当前），单位：毫秒")
	genCmd.Flags().StringVarP(&mnemonic, "mnemonic", "m", "", "助记词")
	genCmd.Flags().StringVarP(&appKey, "app", "a", "", "APP ID")

	//
	genCmd.MarkFlagsMutuallyExclusive("mnemonic", "private")
	genCmd.MarkFlagsMutuallyExclusive("mnemonic", "public")
	genCmd.MarkFlagsRequiredTogether("private", "public")
}

func genRunE(cmd *cobra.Command, args []string) error {
	var public, private []byte
	var err error
	if mnemonic != "" {
		w, err := wallet.NewWalletFromMnemonic(mnemonic)
		if err != nil {
			return err
		}
		public, private = w.GetKeyParis()
	} else {
		public, err = hex.DecodeString(publicKey)
		if err != nil {
			return err
		}

		private, err = hex.DecodeString(privateKey)
		if err != nil {
			return err
		}
	}

	client := auth.NewDefaultApiAuthenticator()
	sig := client.Request(appKey, public, private)
	cmd.Printf("signature is: %s\n", sig)
	return nil
}

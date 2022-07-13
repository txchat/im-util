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
package seed

import (
	"encoding/hex"
	"strings"

	"github.com/33cn/chain33/wallet"
	"github.com/spf13/cobra"
)

// encodingCmd represents the encoding command
var encodingCmd = &cobra.Command{
	Use:     "encoding",
	Short:   "将seed原文，通过密码加密后输出",
	Long:    ``,
	Example: `seed -p '1234qwer' -s '0xFFFFFFF'`,
	Run:     encodingRun,
}

func init() {
	Cmd.AddCommand(encodingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encodingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encodingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	encodingCmd.Flags().StringVarP(&password, "password", "p", "", "password encoded by hex string")
	encodingCmd.Flags().StringVarP(&seed, "seed", "s", "", "seed encoded by hex string")
	encodingCmd.MarkFlagRequired("password")
	encodingCmd.MarkFlagRequired("seed")
}

func encodingRun(cmd *cobra.Command, args []string) {
	seedByte, err := hex.DecodeString(strings.Replace(seed, "0x", "", 1))
	if err != nil {
		cmd.PrintErrf("hex decoding seed failed: %v\n", err)
		return
	}

	encPwd := EncPasswd(password)
	cmd.Printf("encrypted password is: %s\n", hex.EncodeToString(encPwd))

	encSeed, err := wallet.AesgcmEncrypter(encPwd, seedByte)
	if err != nil {
		cmd.PrintErrf("aesGcm encrypt seed failed: %v\n", err)
		return
	}
	cmd.Printf("encrypted seed is: %s\n", hex.EncodeToString(encSeed))
}

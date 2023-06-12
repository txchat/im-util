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
	"github.com/spf13/cobra"
	"github.com/txchat/dtalk/pkg/address"
	"github.com/txchat/dtalk/pkg/auth"
	xcrypt "github.com/txchat/dtalk/pkg/crypt"
	secp256k1_haltingstate "github.com/txchat/dtalk/pkg/crypt/secp256k1-haltingstate"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "校验im鉴权token",
	Long: `通过对指定App下的Token验签，返回签名是否合法。
签名和验签的过程参考：https://github.com/txchat/dtalk/tree/main/pkg/auth`,
	RunE: checkRunE,
}

var (
	token       string
	isCkTimeOut bool
	driver      xcrypt.Encrypt
)

func init() {
	Cmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	checkCmd.Flags().StringVarP(&token, "token", "t", "", "签名数据，原文字符串")
	checkCmd.Flags().StringVarP(&appKey, "app", "a", "dtalk", "APP ID[默认：dtalk]")
	checkCmd.Flags().BoolVarP(&isCkTimeOut, "timeout", "", false, "check timeout enable -t=[false]")

	checkCmd.MarkFlagRequired("token")

	var err error
	driver, err = xcrypt.Load(secp256k1_haltingstate.Name)
	if err != nil {
		panic(err)
	}
}

func checkRunE(cmd *cobra.Command, args []string) error {
	apiRequest, err := auth.NewAPIRequestFromToken(token)
	if err != nil {
		return err
	}
	signatory, err := auth.NewSignatoryFromMetadata(driver, apiRequest.GetMetadata())
	if err != nil {
		return err
	}
	if isMatch, err := signatory.Match(apiRequest.GetSignature(), apiRequest.GetPublicKey()); !isMatch {
		return auth.ErrSignatureInvalid(err)
	}
	if isCkTimeOut && signatory.IsExpire() {
		return auth.ErrSignatureExpired
	}
	uid := address.PublicKeyToAddress(address.NormalVer, apiRequest.GetPublicKey())
	if uid == "" {
		return auth.ErrUIDInvalid
	}
	cmd.Printf("verify success, uid is: %s\n", uid)
	return nil
}

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
package wallet

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/txchat/im-util/protocol/wallet"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "生成钱包",
	Long: `通过指定用户数量，能够生成相应数量的用户wallet。
生成的wallet包含助记词、公私钥、地址信息，这些信息可以通过指定写入文件和分隔符保存在文件中。
同时可以通过读取文件并指定分隔符来将用户wallet从文件中读出。
注意：当指定的数量n大于读取的用户数量时，会自动生成新的用户（填充满不足的用户）；而指定的数量n小于读取的用户数量时则以读取的为准。
每次写入时都是覆盖写入。`,
	Example: `gen -n 1 -d`,
	Run:     genWallet,
}

var (
	writeDir   string
	readDir    string
	newNum     uint32
	readSplit  string
	writeSplit string
	display    bool
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
	genCmd.Flags().StringVarP(&writeDir, "write", "w", "", "[file] write to filed dir")
	genCmd.Flags().StringVarP(&readDir, "read", "r", "", "[file] read from file dir")
	genCmd.Flags().Uint32VarP(&newNum, "numbers", "n", 0, "[must] number of new wallet, default 0")
	genCmd.Flags().StringVarP(&readSplit, "rs", "", ",", "[read must] read split")
	genCmd.Flags().StringVarP(&writeSplit, "ws", "", ",", "[write must] write split")
	genCmd.Flags().BoolVarP(&display, "display", "d", false, "display enable flag")

	genCmd.MarkFlagsRequiredTogether("write", "ws")
	genCmd.MarkFlagsRequiredTogether("read", "rs")
}

func genWallet(cmd *cobra.Command, args []string) {
	var users []*wallet.Wallet
	//读取用户
	if readDir != "" {
		readDriver := wallet.NewFSDriver(readDir, readSplit)
		metadata, err := readDriver.Load()
		if err != nil {
			cmd.PrintErrf("can not load from storage: %v\n", err)
			return
		}
		fmt.Printf("load users:%d\r\n", len(metadata))

		//从文件生成
		fa := wallet.NewFactory(wallet.NewMnemonicCreator(metadata))
		err = fa.Create(runtime.NumCPU())
		if err != nil {
			cmd.PrintErrf("can not create wallet: %v\n", err)
			return
		}
		users = fa.GetRet()
	}

	//生成新用户
	if newNum > uint32(len(users)) {
		fa := wallet.NewFactory(wallet.NewProduceCreator(int(newNum) - len(users)))
		err := fa.Create(runtime.NumCPU())
		if err != nil {
			cmd.PrintErrf("can not create wallet: %v\n", err)
			return
		}
		users = append(users, fa.GetRet()...)
	}

	metadata := make([]*wallet.Metadata, len(users))
	for i, user := range users {
		md, err := wallet.FormatMetadataFromWallet(user)
		if err != nil {
			cmd.PrintErrf("wallet to metadata failed: %v\n", err)
			return
		}
		metadata[i] = md
	}

	//打印
	if display {
		p := wallet.NewPrinter(metadata)
		p.Print()
	}

	//存储用户
	if writeDir != "" {
		saveDriver := wallet.NewFSDriver(writeDir, writeSplit)
		err := saveDriver.Save(metadata)
		if err != nil {
			cmd.PrintErrf("can not save to storage: %v\n", err)
			return
		}
	}
}

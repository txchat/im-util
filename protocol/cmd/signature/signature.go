package signature

import (
	"github.com/spf13/cobra"
	secp256k1_ethereum "github.com/txchat/dtalk/pkg/crypt/secp256k1-ethereum"
	secp256k1_haltingstate "github.com/txchat/dtalk/pkg/crypt/secp256k1-haltingstate"
)

// Cmd represents the signature command
var Cmd = &cobra.Command{
	Use:   "signature",
	Short: "signature是用于标准的签名和验签方法",
	Long: `用于标准的签名和验签方法

默认的签名算法库「secp256k1_haltingstate」
目前支持：
1. 「secp256k1_haltingstate」
2. 「secp256k1_ethereum」
`,
}

type algorithmType int

var (
	algorithm  int
	msg        string
	sig        string
	privateKey string
	publicKey  string

	algorithmName = map[algorithmType]string{
		1: secp256k1_haltingstate.Name,
		2: secp256k1_ethereum.Name,
	}
)

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tokenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tokenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

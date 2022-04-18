package dh

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	walletapi "github.com/txchat/chatcipher"
)

var SessionKey = &cobra.Command{
	Use:     "gen",
	Short:   "generate dh session key",
	Long:    "",
	Example: "  imall dh gen -p ... -P ...",
	Run:     GenerateDHSessionKey,
}

var (
	privateKey string
	publicKey  string
)

func init() {
	SessionKey.Flags().StringVarP(&privateKey, "private", "p", "", "input private key")
	SessionKey.Flags().StringVarP(&publicKey, "public", "P", "", "input public key")
}

func GenerateDHSessionKey(cmd *cobra.Command, args []string) {
	privateKey = strings.Replace(privateKey, "0x", "", 1)
	publicKey = strings.Replace(publicKey, "0x", "", 1)

	sessionKey, err := walletapi.GenerateDHSessionKey(privateKey, publicKey)
	if err != nil {
		fmt.Printf("walletapi.GenerateDHSessionKey failed err:%v\n", err)
		return
	}
	fmt.Printf("gen session key success!:%s\n", hex.EncodeToString(sessionKey))
	return
}

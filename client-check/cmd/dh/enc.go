package dh

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	walletapi "github.com/txchat/chatcipher"
	"github.com/txchat/chatcipher/service/dh"
)

var Encrypt = &cobra.Command{
	Use:     "enc",
	Short:   "encrypt by dh session key",
	Long:    "",
	Example: "  imall dh enc -k ... -d ...",
	Run:     DoEncrypt,
}

var (
	data       string
	sessionKey string
)

func init() {
	Encrypt.Flags().StringVarP(&data, "data", "d", "", "input source data")
	Encrypt.Flags().StringVarP(&sessionKey, "key", "k", "", "input session key")
	Encrypt.Flags().StringVarP(&privateKey, "private", "p", "", "input private key")
	Encrypt.Flags().StringVarP(&publicKey, "public", "P", "", "input public key")
}

func DoEncrypt(cmd *cobra.Command, args []string) {
	privateKey = strings.Replace(privateKey, "0x", "", 1)
	publicKey = strings.Replace(publicKey, "0x", "", 1)
	sessionKey = strings.Replace(sessionKey, "0x", "", 1)
	d, err := hex.DecodeString(strings.Replace(data, "0x", "", 1))
	if err != nil {
		fmt.Printf("decode srouce data failed err:%v\n", err)
		return
	}
	if privateKey != "" && publicKey != "" {
		rlt, err := walletapi.EncryptWithDHKeyPair(privateKey, publicKey, d)
		if err != nil {
			fmt.Printf("walletapi.GenerateDHSessionKey failed err:%v\n", err)
			return
		}
		fmt.Printf("encrypt by dh key pair success!:%s\n", hex.EncodeToString(rlt))
	} else {
		rlt, err := dh.EncryptSymmetric(sessionKey, d)
		if err != nil {
			fmt.Printf("dh.EncryptSymmetric failed err:%v\n", err)
			return
		}
		fmt.Printf("encrypt by dh session key success!:%s\n", hex.EncodeToString(rlt))
	}
	return
}

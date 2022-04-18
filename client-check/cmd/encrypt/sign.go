package encrypt

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/txchat/im-util/client-check/encrypt"
	//secpEth "github.com/txchat/im-util/client-check/encrypt/secp256k1-ethereum"
	secpBty "github.com/txchat/im-util/client-check/encrypt/secp256k1-haltingstate"
)

var SignCmd = &cobra.Command{
	Use:     "sign",
	Short:   "signed by encryption algorithm",
	Long:    "signed by encryption algorithm",
	Example: "  imall encrypt sign -p [private key] -m [content] -t eth",
	Run:     sign,
}

var (
	privateKey string
	msg        string
	tp         string
)

func init() {
	SignCmd.Flags().StringVarP(&privateKey, "private", "p", "", "private key")
	SignCmd.Flags().StringVarP(&msg, "msg", "m", "", "msg content")
	SignCmd.Flags().StringVarP(&tp, "type", "t", "", "encrypt type: [eth] [bty]")
}

func sign(cmd *cobra.Command, args []string) {
	encType := ""
	if tp == "eth" {
		//encType = secpEth.Name
	}
	if tp == "bty" {
		encType = secpBty.Name
	}
	encrypt, err := encrypt.Load(encType)
	if err != nil {
		fmt.Printf("load sign exec failed err:%v\n type:%v\n", err, encType)
		return
	}
	privKey, err := hex.DecodeString(strings.Replace(privateKey, "0x", "", 1))
	if err != nil {
		fmt.Printf("sign failed err:%v\n", err)
		return
	}
	msg256 := sha256.Sum256([]byte(msg))
	fmt.Println(encrypt.Sign(msg256[:], privKey))
}

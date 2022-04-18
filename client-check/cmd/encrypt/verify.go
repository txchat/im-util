package encrypt

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	//secpEth "github.com/txchat/im-util/client-check/encrypt/secp256k1-ethereum"
	secpBty "github.com/txchat/im-util/client-check/encrypt/secp256k1-haltingstate"
	"strings"

	"github.com/spf13/cobra"
	"github.com/txchat/im-util/client-check/encrypt"
)

var VerifyCmd = &cobra.Command{
	Use:     "verify",
	Short:   "verified by encryption algorithm",
	Long:    "verified by encryption algorithm",
	Example: "  imall encrypt verify -p [private key] -m [content] -t eth",
	Run:     verify,
}

var (
	publicKey string
	sig       string
)

func init() {
	VerifyCmd.Flags().StringVarP(&publicKey, "public", "p", "", "public key")
	VerifyCmd.Flags().StringVarP(&sig, "sig", "s", "", "sig content")
	VerifyCmd.Flags().StringVarP(&msg, "msg", "m", "", "msg content")
	VerifyCmd.Flags().StringVarP(&tp, "type", "t", "", "encrypt type: [eth] [bty]")
}

func verify(cmd *cobra.Command, args []string) {
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
	fmt.Println("input: sig = ", sig)
	fmt.Println("input: public key = ", publicKey)
	fmt.Println("input: msg = ", msg)
	pubKey, err := hex.DecodeString(strings.Replace(publicKey, "0x", "", 1))
	if err != nil {
		fmt.Printf("sign failed err:%v\n", err)
		return
	}
	sigData, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		fmt.Printf("decode sig data failed err:%v\n", err)
		return
	}
	msg256 := sha256.Sum256([]byte(msg))
	fmt.Println(encrypt.Verify(msg256[:], sigData, pubKey))
}

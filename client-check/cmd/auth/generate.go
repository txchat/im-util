package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/spf13/cobra"
	walletapi "github.com/txchat/chatcipher"
	"strings"
)

var GenCmd = &cobra.Command{
	Use:     "gen",
	Short:   "generate auth token",
	Long:    "generate auth token",
	Example: " imall gen token ",
	Run:     doGenCmd,
}

var (
	privateKey string
	publicKey  string
	mnemonic   string
	timestamp  int64
	randomStr  string
)

func init() {
	GenCmd.Flags().StringVarP(&privateKey, "private", "p", "", "private key encoded by hex")
	GenCmd.Flags().Int64VarP(&timestamp, "timestamp", "t", 0, "timestamp")
	GenCmd.Flags().StringVarP(&randomStr, "random", "r", "", "random string")
	GenCmd.Flags().StringVarP(&publicKey, "public", "P", "", "public key encoded by Hex")
	GenCmd.Flags().StringVarP(&mnemonic, "mnemonic", "m", "", "mnemonic string")
}

func doGenCmd(cmd *cobra.Command, args []string) {
	var privKey, pubKey []byte
	var err error
	if privateKey == "" || publicKey == "" {
		privKey, pubKey, err = mnemonicToKeyPairs(mnemonic)
		if err != nil {
			fmt.Printf("generate token failed err:%v\n", err)
			return
		}
	} else {
		privKey, err = hex.DecodeString(strings.Replace(privateKey, "0x", "", 1))
		if err != nil {
			fmt.Printf("generate token failed err:%v\n", err)
			return
		}
		pubKey, err = hex.DecodeString(strings.Replace(publicKey, "0x", "", 1))
		if err != nil {
			fmt.Printf("generate token failed err:%v\n", err)
			return
		}
	}
	message := fmt.Sprintf("%d*%s", timestamp, randomStr)
	msg256 := sha256.Sum256([]byte(message))
	sig := base64.StdEncoding.EncodeToString(walletapi.ChatSign(msg256[:], privKey))
	fmt.Println(fmt.Sprintf("%s#%s#%s", sig, message, hex.EncodeToString(pubKey)))
}

func mnemonicToKeyPairs(mnemonic string) ([]byte, []byte, error) {
	wallet, err := walletapi.NewWalletFromMnemonic_v2("BTC", mnemonic)
	if err != nil {
		return nil, nil, err
	}

	privKey, err := wallet.NewKeyPriv(0)
	if err != nil {
		return nil, nil, err
	}

	pubKey, err := wallet.NewKeyPub(0)
	if err != nil {
		return nil, nil, err
	}
	return privKey, pubKey, nil
}

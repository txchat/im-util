package seed

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	walletapi "github.com/txchat/chatcipher"
)

var seedEnc = &cobra.Command{
	Use:     "encode",
	Short:   "encode",
	Long:    "get encoded seed by password",
	Example: "encode -p '1234qwer' -s '0xFFFFFFF......FFF' ",
	Run:     doSeedEnc,
}

func init() {
	seedEnc.Flags().StringVarP(&password, "password", "p", "", "password encoded by hex string")
	seedEnc.Flags().StringVarP(&seed, "seed", "s", "", "seed encoded by hex string")
}

func doSeedEnc(cmd *cobra.Command, args []string) {
	seedData, err := hex.DecodeString(strings.Replace(seed, "0x", "", 1))
	if err != nil {
		fmt.Printf("get seed failed err:%v\n", err)
		return
	}

	pwdData := walletapi.EncPasswd(password)
	fmt.Printf("encPasswd reslut:%s\n", hex.EncodeToString(pwdData))
	data, err := walletapi.SeedEncKey(pwdData, seedData)
	if err != nil {
		fmt.Printf("get seed failed err:%v\n", err)
		return
	}
	fmt.Printf("get seed success!:%s\n", hex.EncodeToString(data))
	return
}

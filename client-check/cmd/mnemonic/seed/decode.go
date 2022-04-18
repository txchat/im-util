package seed

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	walletapi "github.com/txchat/chatcipher"
)

var (
	password string
	seed     string
)

func init() {
	seedDec.Flags().StringVarP(&password, "password", "p", "", "password encoded by hex string")
	seedDec.Flags().StringVarP(&seed, "seed", "s", "", "seed encoded by hex string")
}

var seedDec = &cobra.Command{
	Use:     "decode",
	Short:   "decode",
	Long:    "get decoded seed by password",
	Example: "decode -d fffff",
	Run:     doSeedDec,
}

func doSeedDec(cmd *cobra.Command, args []string) {
	seedData, err := hex.DecodeString(strings.Replace(seed, "0x", "", 1))
	if err != nil {
		fmt.Printf("get seed failed err:%v\n", err)
		return
	}
	pwdData := walletapi.EncPasswd(password)
	fmt.Printf("encPasswd reslut:%s\n", hex.EncodeToString(pwdData))
	ret, err := walletapi.SeedDecKey(pwdData, seedData)
	if err != nil {
		fmt.Printf("SeedDecKey failed err:%v\n", err)
		return
	}
	fmt.Printf("get seed success!:%s\n", hex.EncodeToString(ret))
	return
}

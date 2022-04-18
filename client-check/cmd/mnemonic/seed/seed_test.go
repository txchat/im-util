package seed

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"

	walletapi "github.com/txchat/chatcipher"
)

func Test_decodeSeedFunc(t *testing.T) {
	password := "1234qwer"
	seed := "0xdd19aa8a5adfb09f9b0bb654c59efde6ad9cd76c9326b4dccb731c726b01175e957a950c4e36d004484a9857da66f536092569fb421b4d615e57fc45a2"
	seedData, err := hex.DecodeString(strings.Replace(seed, "0x", "", 1))
	if err != nil {
		t.Error(fmt.Printf("get seed failed err:%v\n", err))
		return
	}
	//pwdData, err := hex.DecodeString(password)
	//if err != nil {
	//	fmt.Printf("get password failed err:%v\n", err)
	//	return
	//}
	now := time.Now()
	pwdData := walletapi.EncPasswd(password)
	time.Since(now)
	ret, err := walletapi.SeedDecKey(pwdData, seedData)
	if err != nil {
		t.Error(fmt.Printf("SeedDecKey failed err:%v\n", err))
		return
	}
	t.Log(fmt.Printf("get seed success!:%s\n", hex.EncodeToString(ret)))
	return
}

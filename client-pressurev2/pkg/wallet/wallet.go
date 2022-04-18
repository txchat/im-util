package wallet

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/txchat/chatcipher/crypto/secp256k1"
)

type Wallet struct {
	Address string
	PubKey  []byte
	PrivKey []byte
	Mem     string
}

func GetToken(privateKey, publicKey []byte) string {
	datatime := time.Now().UnixNano() / 1000
	msg := []byte(strconv.FormatInt(datatime, 10) + "*" + randomString(10))
	tmsg := sha256.Sum256(msg)
	sig := chatSign(tmsg[:], privateKey)
	sig64 := base64.StdEncoding.EncodeToString(sig)

	pubKey := hex.EncodeToString(publicKey)

	return sig64 + "#" + string(msg) + "#" + pubKey
}

// chatSign 签名
func chatSign(msg, privateKey []byte) []byte {
	res, err := secp256k1.Sign(msg, privateKey)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func randomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

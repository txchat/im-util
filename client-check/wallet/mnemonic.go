package wallet

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/33cn/chain33/common/address"
	cipher "github.com/txchat/chatcipher"
	"github.com/txchat/chatcipher/crypto/secp256k1"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Wallet struct {
	wallet     *cipher.HDWallet
	privateKey []byte
	publicKey  []byte
	mnemonic   string
}

func NewWallet(mnemonic string) (*Wallet, error) {
	wallet, err := cipher.NewWalletFromMnemonic_v2("BTY", mnemonic)
	if err != nil {
		return nil, err
	}

	privKey, err := wallet.NewKeyPriv(0)
	if err != nil {
		return nil, err
	}

	pubKey, err := wallet.NewKeyPub(0)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		wallet:     wallet,
		privateKey: privKey,
		publicKey:  pubKey,
		mnemonic:   mnemonic,
	}, nil
}

func (c *Wallet) GetPrivateKey() []byte {
	return c.privateKey
}

func (c *Wallet) GetPublicKey() []byte {
	return c.publicKey
}

func (c *Wallet) GetMnemonic() string {
	return c.mnemonic
}

func (c *Wallet) GetAddress() string {
	return address.PubKeyToAddress(c.publicKey).String()
}

func (c *Wallet) JoinString(split string) string {
	return strings.Join([]string{c.mnemonic, hex.EncodeToString(c.privateKey), hex.EncodeToString(c.publicKey), c.GetAddress()}, split)
}

func (c *Wallet) GetToken() string {
	datatime := time.Now().UnixNano() / 1000000
	msg := []byte(strconv.FormatInt(datatime, 10) + "*" + randomString(10))
	tmsg := sha256.Sum256(msg)
	sig := chatSign(tmsg[:], c.privateKey)
	sig64 := base64.StdEncoding.EncodeToString(sig)

	pubKey := hex.EncodeToString(c.publicKey)

	return sig64 + "#" + string(msg) + "#" + pubKey
}

// chatSign 签名
func chatSign(msg, privateKey []byte) []byte {
	res, _ := secp256k1.Sign(msg, privateKey)
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

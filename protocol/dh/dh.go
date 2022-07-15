package dh

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/txchat/im-util/protocol/dh/ecies"
)

const (
	aesNonceLength = 12
	//aesKeyLength   = 32 // in bytes

	// Shared secret key length
	sskLen = 16
)

//GenerateDHSessionKey : 根据本端私钥和对端公钥生成ECDH会话密钥
func GenerateDHSessionKey(privateKey, publicKey string) ([]byte, error) {
	ecdsaPrivateKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}
	eciesPrivateKey := ecies.ImportECDSA(ecdsaPrivateKey)

	publicKeyData, err := hex.DecodeString(publicKey)
	if err != nil {
		return nil, err
	}
	ecdsaPublicKey, err := crypto.DecompressPubkey(publicKeyData)
	if err != nil {
		return nil, err
	}
	eciesPublicKey := ecies.ImportECDSAPublic(ecdsaPublicKey)

	return eciesPrivateKey.GenerateShared(eciesPublicKey, sskLen, sskLen)
}

//EncryptWithDHKeyPair : 根据用户的私钥和对端的公钥生成ecdh密钥并进行对称加密
func EncryptWithDHKeyPair(privateKey, publicKey string, plaintext []byte) ([]byte, error) {
	key, err := GenerateDHSessionKey(privateKey, publicKey)
	if err != nil {
		return nil, err
	}

	return encryptSymmetric(key, plaintext)
}

//DecryptWithDHKeyPair : 根据用户的私钥和对端的公钥生成ecdh密钥并进行对称解密
func DecryptWithDHKeyPair(privateKey, publicKey string, cyphertext []byte) ([]byte, error) {
	key, err := GenerateDHSessionKey(privateKey, publicKey)
	if err != nil {
		return nil, err
	}

	return decryptSymmetric(key, cyphertext)
}

func EncryptSymmetric(key string, plaintext []byte) ([]byte, error) {
	keySlice, err := hex.DecodeString(strings.ReplaceAll(key, "0x", key))
	if nil != err {
		return nil, err
	}

	return encryptSymmetric(keySlice, plaintext)
}

func DecryptSymmetric(key string, cyphertext []byte) ([]byte, error) {
	keySlice, err := hex.DecodeString(strings.ReplaceAll(key, "0x", key))
	if nil != err {
		return nil, err
	}

	return decryptSymmetric(keySlice, cyphertext)
}

//encryptSymmetric : 对称加密
func encryptSymmetric(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, aesNonceLength)
	if _, err = rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("cannot read random bytes for nonce: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	encrypted, err := aesgcm.Seal(nil, nonce, plaintext, nil), nil
	if err != nil {
		return nil, err
	}

	return append(encrypted, nonce...), nil
}

//decryptSymmetric : 对称解密
func decryptSymmetric(key []byte, cyphertext []byte) ([]byte, error) {
	// symmetric messages are expected to contain the 12-byte nonce at the end of the payload
	if len(cyphertext) < aesNonceLength {
		return nil, fmt.Errorf("missing salt or invalid payload in symmetric message")
	}
	salt := cyphertext[len(cyphertext)-aesNonceLength:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	decrypted, err := aesgcm.Open(nil, salt, cyphertext[:len(cyphertext)-aesNonceLength], nil)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

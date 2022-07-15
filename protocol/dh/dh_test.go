package dh

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mne1        = "担 将 甜 机 打 宴 屋 藏 罚 燕 欧 市 喂 障 颜"
	privateKey1 = "2700c94e67119e2c523caf8aed574e5bf3d633b261df4d451128f8f284007f4b"
	publicKey1  = "037fb52a7e8c89151fc6460513f585c8b82cd70a363003bc33a341003fc3a412b2"
	btcAddr1    = "17j3vDeEqqWbd3kePChgeTWn9Rg4JiUsyM"

	mne2        = "疏 跨 糊 挥 刮 害 若 皱 往 候 姚 遍 乏 代 墨"
	privateKey2 = "42050e738eedd791d2d827b443ab3aee2ff893add73f0c0187ce973ba46f9603"
	publicKey2  = "02fbe91249dc684114381c1cc9e0d8edb92f962175318b1d08aea32c99fcd246a9"
	btcAddr2    = "1HvCD2uZoiGC7oUygB88k5uqeMn2wQtZx5"

	msg             = "hello world"
	encryptedSource = "975af04e7a7a1fdd290175064b48ffca47ae36c72678bd064e67bf19b91b6790087012173e427f"
	sskSource       = "f6a9bb98468506b517d9d9e9ad31b699f9afb6f05599760a68f1a39758f529e1"
)

func TestGenerateDHSessionKey(t *testing.T) {
	ssk, err := GenerateDHSessionKey(privateKey1, publicKey2)
	assert.Nil(t, err)
	assert.Equal(t, hex.EncodeToString(ssk), sskSource)
}

func TestEncryptWithDHKeyPair(t *testing.T) {
	encryptedData, err := EncryptWithDHKeyPair(privateKey1, publicKey2, []byte(msg))
	assert.Nil(t, err)
	assert.NotNil(t, encryptedData)
	decryptedData, err := DecryptWithDHKeyPair(privateKey2, publicKey1, encryptedData)
	assert.Nil(t, err)
	assert.Equal(t, decryptedData, []byte(msg))
}

func TestDecryptWithDHKeyPair(t *testing.T) {
	encryptedSourceData, err := hex.DecodeString(encryptedSource)
	assert.Nil(t, err)
	decryptedData, err := DecryptWithDHKeyPair(privateKey2, publicKey1, encryptedSourceData)
	assert.Nil(t, err)
	assert.Equal(t, decryptedData, []byte(msg))
}

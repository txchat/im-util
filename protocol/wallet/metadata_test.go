package wallet

import (
	"encoding/hex"
	"testing"

	_ "github.com/33cn/chain33/system/address/btc"
	_ "github.com/33cn/chain33/system/address/eth"
	"github.com/stretchr/testify/assert"
)

var (
	mne        = "担 将 甜 机 打 宴 屋 藏 罚 燕 欧 市 喂 障 颜"
	privateKey = "2700c94e67119e2c523caf8aed574e5bf3d633b261df4d451128f8f284007f4b"
	publicKey  = "037fb52a7e8c89151fc6460513f585c8b82cd70a363003bc33a341003fc3a412b2"
	btcAddr    = "17j3vDeEqqWbd3kePChgeTWn9Rg4JiUsyM"
	ethAddr    = "0x98182d1ced67e003b57cf1c44fd28ab3c7b33752"
)

func TestFormatMetadataFromWallet(t *testing.T) {
	private, err := hex.DecodeString(privateKey)
	assert.Nil(t, err)
	public, err := hex.DecodeString(publicKey)
	assert.Nil(t, err)
	w, err := NewWalletFromMetadata(&Metadata{
		privateKey: private,
		publicKey:  public,
		mnemonic:   mne,
	})
	assert.Nil(t, err)
	//bty address
	md, err := FormatMetadataFromWallet(0, w)
	assert.Nil(t, err)
	assert.Equal(t, btcAddr, md.GetAddress())

	//eth address
	md2, err := FormatMetadataFromWallet(2, w)
	assert.Nil(t, err)
	assert.Equal(t, ethAddr, md2.GetAddress())
}

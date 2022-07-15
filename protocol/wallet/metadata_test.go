package wallet

import (
	"encoding/hex"
	"testing"

	_ "github.com/33cn/chain33/system/address/btc"
	_ "github.com/33cn/chain33/system/address/eth"
	"github.com/stretchr/testify/assert"
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

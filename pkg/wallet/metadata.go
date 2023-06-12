package wallet

import (
	"fmt"

	"github.com/33cn/chain33/common/address"
)

type Metadata struct {
	privateKey []byte
	publicKey  []byte
	mnemonic   string
	address    string
}

func FormatMetadataFromWallet(addrType int32, w *Wallet) (*Metadata, error) {
	addr := address.PubKeyToAddr(addrType, w.publicKey)
	if addr == "" {
		return nil, fmt.Errorf("address type[%d] generate empty address", addrType)
	}
	return &Metadata{
		privateKey: w.privateKey,
		publicKey:  w.publicKey,
		mnemonic:   w.mnemonic,
		address:    addr,
	}, nil
}

func (m *Metadata) GetPrivateKey() []byte {
	return m.privateKey
}

func (m *Metadata) GetPublicKey() []byte {
	return m.publicKey
}

func (m *Metadata) GetMnemonic() string {
	return m.mnemonic
}

func (m *Metadata) GetAddress() string {
	return m.address
}

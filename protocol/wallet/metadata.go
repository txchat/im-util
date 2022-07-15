package wallet

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

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

func FormatMetadata(row, split string) (*Metadata, error) {
	var privateKey []byte
	var publicKey []byte
	mne := row
	addr := ""

	if split != "" {
		var err error
		items := strings.Split(row, split)
		if len(items) != 4 {
			return nil, errors.New("item number is not 4")
		}
		mne = items[0]
		addr = items[3]
		privateKey, err = hex.DecodeString(items[1])
		if err != nil {
			return nil, err
		}
		publicKey, err = hex.DecodeString(items[2])
		if err != nil {
			return nil, err
		}
	}
	return &Metadata{
		privateKey: privateKey,
		publicKey:  publicKey,
		mnemonic:   mne,
		address:    addr,
	}, nil
}

func (m *Metadata) Convert(split string) string {
	if split != "" {
		return strings.Join([]string{m.mnemonic, hex.EncodeToString(m.privateKey), hex.EncodeToString(m.publicKey), m.address}, split)
	}
	return m.mnemonic
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

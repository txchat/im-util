package wallet

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/33cn/chain33/wallet/bipwallet"
)

type Metadata struct {
	privateKey []byte
	publicKey  []byte
	mnemonic   string
	address    string
}

func FormatMetadataFromWallet(w *Wallet) (*Metadata, error) {
	address, err := bipwallet.PubToAddress(w.publicKey)
	if err != nil {
		return nil, err
	}
	return &Metadata{
		privateKey: w.privateKey,
		publicKey:  w.publicKey,
		mnemonic:   w.mnemonic,
		address:    address,
	}, nil
}

func FormatMetadata(row, split string) (*Metadata, error) {
	var privateKey []byte
	var publicKey []byte
	mne := row
	address := ""

	if split != "" {
		var err error
		items := strings.Split(row, split)
		if len(items) != 4 {
			return nil, errors.New("item number is not 4")
		}
		mne = items[0]
		address = items[3]
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
		address:    address,
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

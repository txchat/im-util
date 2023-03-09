package wallet

import (
	"encoding/hex"
	"errors"
	"strings"
)

type Formatter interface {
	FromRow(row string) (*Metadata, error)
	ToRow(metadata *Metadata) string
}

type SplitFormatter struct {
	split string
}

func NewSplitFormatter(split string) *SplitFormatter {
	return &SplitFormatter{
		split: split,
	}
}

func (sf *SplitFormatter) FromRow(row string) (*Metadata, error) {
	var privateKey []byte
	var publicKey []byte
	mne := row
	addr := ""

	if sf.split != "" {
		var err error
		items := strings.Split(row, sf.split)
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

func (sf *SplitFormatter) ToRow(m *Metadata) string {
	if sf.split != "" {
		return strings.Join([]string{m.mnemonic, hex.EncodeToString(m.privateKey), hex.EncodeToString(m.publicKey), m.address}, sf.split)
	}
	return m.mnemonic
}

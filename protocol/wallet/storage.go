package wallet

import (
	"bufio"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"os"
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

type FSDriver struct {
	split string
	uri   string

	rows []string
}

func NewFSDriver(uri, split string) *FSDriver {
	return &FSDriver{
		split: split,
		uri:   uri,
	}
}

func (d *FSDriver) Load() ([]*Metadata, error) {
	// load
	var err error
	d.rows, err = d.readLines(d.uri)
	if err != nil {
		return nil, err
	}

	ret := make([]*Metadata, len(d.rows))
	for i, row := range d.rows {
		ret[i], err = FormatMetadata(row, d.split)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (d *FSDriver) Save(metadata []*Metadata) error {
	f, err := os.OpenFile(d.uri, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rb := bufio.NewWriter(f)
	defer rb.Flush()

	d.rows = make([]string, len(metadata))
	for i, md := range metadata {
		item := md.Convert(d.split)
		_, err := rb.WriteString(item + "\n")
		if err != nil {
			return err
		}
		d.rows[i] = item
	}
	return nil
}

func (d *FSDriver) readLines(uri string) ([]string, error) {
	f, err := os.OpenFile(uri, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	list := strings.Split(string(b), "\n")
	if len(list) > 0 && list[len(list)-1] == "" {
		return list[:len(list)-1], nil
	}
	return list, nil
}

package reader

import (
	"errors"
	"github.com/txchat/im-util/protocol/wallet"
)

func LoadMetadata(readPath, readSplit string) ([]*wallet.Metadata, error) {
	if readPath == "" || readSplit == "" {
		return nil, errors.New("readPath or readSplit can not empty")
	}
	readDriver := wallet.NewFSDriver(readPath, readSplit)
	return readDriver.Load()
}

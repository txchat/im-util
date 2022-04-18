package filehelper

import (
	"encoding/hex"
	"github.com/txchat/im-util/client-pressurev2/pkg/wallet"
	"io/ioutil"
	"os"
	"strings"
)

type FileReader struct {
	lines []string
}

func NewFileReader() *FileReader {
	return &FileReader{lines: make([]string, 0)}
}

func (f *FileReader) ReadFile(filepath string) error {
	fd, err := os.OpenFile(filepath, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer fd.Close()
	b, err := ioutil.ReadAll(fd)
	if err != nil {
		return err
	}
	f.lines = strings.Split(string(b), "\n")
	return nil
}

func (f *FileReader) GetAllLines() []string {
	return f.lines
}

func (f *FileReader) GetUserWallet(num int) []*wallet.Wallet {
	var store = make([]*wallet.Wallet, 0)
	count := 0
	for _, row := range f.lines {
		if count >= num {
			break
		}
		// 拆分 DB.txt
		items := strings.Split(row, ",")
		if len(items) != 4 {
			continue
		}
		// 获得公钥私钥后生成 token
		// mnemonic := items[0]
		priKey, _ := hex.DecodeString(items[1])
		pubKey, _ := hex.DecodeString(items[2])
		addr := items[3]
		store = append(store, &wallet.Wallet{
			Address: addr,
			PubKey:  pubKey,
			PrivKey: priKey,
			Mem:     "",
		})
		count++
	}
	return store
}

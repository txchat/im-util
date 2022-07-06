package wallet

/*
	将用户助记词等元数据以特定的格式输出到屏幕
*/
import (
	"encoding/hex"
	"fmt"
)

type Printer struct {
	mode int
	data []*Metadata
}

func NewPrinter(data []*Metadata) *Printer {
	return &Printer{
		mode: 0,
		data: data,
	}
}

func (p *Printer) Print() {
	for _, d := range p.data {
		fmt.Printf("助记词：%s;私钥：%s;公钥：%s;地址：%s\r\n",
			d.mnemonic,
			hex.EncodeToString(d.privateKey),
			hex.EncodeToString(d.publicKey),
			d.address)
	}
}

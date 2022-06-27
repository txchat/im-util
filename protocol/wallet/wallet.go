package wallet

import (
	"github.com/33cn/chain33/types"
	"github.com/33cn/chain33/wallet/bipwallet"
	"sync"
)

type Wallet struct {
	wallet     *bipwallet.HDWallet
	privateKey []byte
	publicKey  []byte
	mnemonic   string
}

func CreateNewWallet() (*Wallet, error) {
	//创建助记词
	mne, err := bipwallet.NewMnemonicString(1, 160)
	if err != nil {
		return nil, err
	}
	//创建钱包
	wallet, err := bipwallet.NewWalletFromMnemonic(bipwallet.TypeBty, uint32(types.SECP256K1), mne)
	if err != nil {
		return nil, err
	}
	priv, pub, err := wallet.NewKeyPair(0)
	if err != nil {
		return nil, err
	}
	w := Wallet{
		wallet:     wallet,
		privateKey: priv,
		publicKey:  pub,
		mnemonic:   mne,
	}
	return &w, nil
}

func NewWalletFromMetadata(data *Metadata) (*Wallet, error) {
	//创建钱包
	wallet, err := bipwallet.NewWalletFromMnemonic(bipwallet.TypeBty, uint32(types.SECP256K1), data.mnemonic)
	if err != nil {
		return nil, err
	}
	if data.privateKey == nil || data.publicKey == nil {
		var err error
		data.privateKey, data.publicKey, err = wallet.NewKeyPair(0)
		if err != nil {
			return nil, err
		}
	}
	w := Wallet{
		wallet:     wallet,
		privateKey: data.privateKey,
		publicKey:  data.publicKey,
		mnemonic:   data.mnemonic,
	}
	return &w, nil
}

func NewWalletFromMnemonic(mnemonic string) (*Wallet, error) {
	return NewWalletFromMetadata(&Metadata{
		mnemonic: mnemonic,
	})
}

// GetKeyParis publicKey, privateKey
func (w *Wallet) GetKeyParis() ([]byte, []byte) {
	return w.publicKey, w.privateKey
}

type Creator interface {
	Length() int
	Foreach(start, end int) Iterator
	NewWallet(v interface{}) (*Wallet, error)
}

type Iterator interface {
	HasNext() bool
	Next() (int, interface{})
}

type MnemonicCreatorIterator struct {
	data  []*Metadata
	index int
}

func (i *MnemonicCreatorIterator) HasNext() bool {
	return i.index < len(i.data)
}

func (i *MnemonicCreatorIterator) Next() (index int, v interface{}) {
	index = i.index
	v = i.data[i.index]
	i.index++
	return
}

type MnemonicCreator struct {
	mds []*Metadata
}

func NewMnemonicCreator(mds []*Metadata) *MnemonicCreator {
	return &MnemonicCreator{
		mds: mds,
	}
}

func (mc *MnemonicCreator) Length() int {
	return len(mc.mds)
}

func (mc *MnemonicCreator) Foreach(start, end int) Iterator {
	return &MnemonicCreatorIterator{
		data:  mc.mds[start:end],
		index: 0,
	}
}

func (mc *MnemonicCreator) NewWallet(v interface{}) (*Wallet, error) {
	return NewWalletFromMetadata(v.(*Metadata))
}

//
type ProduceCreatorIterator struct {
	len   int
	start int
	index int
}

func (i *ProduceCreatorIterator) HasNext() bool {
	return i.index < i.len
}

func (i *ProduceCreatorIterator) Next() (index int, v interface{}) {
	index = i.index
	v = i.start + i.index
	i.index++
	return
}

type ProduceCreator struct {
	number int
}

func NewProduceCreator(number int) *ProduceCreator {
	return &ProduceCreator{
		number: number,
	}
}

func (mc *ProduceCreator) Length() int {
	return mc.number
}

func (mc *ProduceCreator) Foreach(start, end int) Iterator {
	return &ProduceCreatorIterator{
		len:   end - start,
		start: start,
		index: 0,
	}
}

func (mc *ProduceCreator) NewWallet(v interface{}) (*Wallet, error) {
	return CreateNewWallet()
}

//
type Factory struct {
	creator Creator
	ret     []*Wallet
}

func NewFactory(creator Creator) *Factory {
	return &Factory{
		creator: creator,
		ret:     make([]*Wallet, creator.Length()),
	}
}

func (f *Factory) GetRet() []*Wallet {
	return f.ret
}

func (f *Factory) Create(cpus int) error {
	wg := sync.WaitGroup{}
	length := f.creator.Length()

	cnt := length / cpus
	last := length % cpus
	for i := 0; i < cpus; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			iter := f.creator.Foreach(index*cnt, (index+1)*cnt)
			for iter.HasNext() {
				var err error
				j, v := iter.Next()
				f.ret[index+j], err = f.creator.NewWallet(v)
				if err != nil {
					panic(err)
				}
			}
		}(i)
	}
	lastStart := length - last
	iter := f.creator.Foreach(lastStart, length)
	for iter.HasNext() {
		var err error
		i, v := iter.Next()
		f.ret[lastStart+i], err = f.creator.NewWallet(v)
		if err != nil {
			return err
		}
	}
	wg.Wait()
	return nil
}

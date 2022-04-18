package address

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	cipher "github.com/txchat/chatcipher"

	"github.com/txchat/im-util/client-check/model"
	"github.com/txchat/im-util/client-check/wallet"
)

var GenCmd = &cobra.Command{
	Use:     "gen",
	Short:   "generate BTY address",
	Long:    "generate BTY address",
	Example: "gen -n 1 --hm=true -d true",
	Run:     addressGen,
}

var (
	writeDir      string
	readDir       string
	newNum        uint32
	readSplit     string
	writeSplit    string
	humanReadable bool
	display       bool
)

func init() {
	GenCmd.Flags().StringVarP(&writeDir, "write", "w", "", "[file] write to filed dir")
	GenCmd.Flags().StringVarP(&readDir, "read", "r", "", "[file] read from file dir")
	GenCmd.Flags().Uint32VarP(&newNum, "numbers", "n", 0, "[must] number of new wallet, default 0")
	GenCmd.Flags().StringVarP(&readSplit, "rs", "", "", "[read must] read split")
	GenCmd.Flags().StringVarP(&writeSplit, "ws", "", "", "[write must] write split")
	GenCmd.Flags().BoolVarP(&humanReadable, "hm", "", false, "human readable enable flag")
	GenCmd.Flags().BoolVarP(&display, "display", "d", false, "display enable flag")
}

func addressGen(cmd *cobra.Command, args []string) {
	var wallets = make([]*wallet.Wallet, 0)
	//读取用户
	if readDir != "" {
		mne, err := readLines(readDir)
		if err != nil {
			log.Error().Err(err).Msg("can not read files")
			return
		}
		fmt.Printf("linse:%d\r\n", len(mne))
		wallets, err = loadClients(readSplit, mne)
		if err != nil {
			log.Error().Err(err).Msg("can not load clients")
			return
		}
	}
	//生成用户
	if newNum > 0 {
		mne, err := createMnemonic(newNum)
		if err != nil {
			log.Error().Err(err).Msg("can not create mnemonics")
			return
		}
		list, err := loadClients("", mne)
		if err != nil {
			log.Error().Err(err).Msg("can not load clients")
			return
		}
		wallets = append(wallets, list...)
	}

	if display {
		for _, w := range wallets {
			if humanReadable {
				fmt.Printf("助记词：%s;私钥：%s;公钥：%s;地址：%s\r\n",
					w.GetMnemonic(), hex.EncodeToString(w.GetPrivateKey()), hex.EncodeToString(w.GetPublicKey()), w.GetAddress())
			} else {
				fmt.Printf("%s\r\n", w.GetMnemonic())
			}
		}
	}

	//存储用户
	if writeDir != "" {
		f, err := os.OpenFile(writeDir, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		rb := bufio.NewWriter(f)
		defer rb.Flush()
		for _, w := range wallets {
			item := ""
			if writeSplit != "" {
				item = w.JoinString(writeSplit)
			} else {
				item = w.GetMnemonic()
			}
			_, err := rb.WriteString(item + "\n")
			if err != nil {
				log.Error().Err(err).Msg("write files failed")
				return
			}
		}
	}
}

func readLines(uri string) ([]string, error) {
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

func createMnemonic(num uint32) ([]string, error) {
	var mnemonics = make([]string, num)
	for i := 0; i < int(num); i++ {
		//创建助记词
		mne, err := cipher.NewMnemonicString(1, 160)
		if err != nil {
			return nil, err
		}
		mnemonics[i] = mne
	}
	return mnemonics, nil
}

func loadClients(split string, rows []string) ([]*wallet.Wallet, error) {
	m := sync.Mutex{}
	wg := sync.WaitGroup{}
	var wallets = make([]*wallet.Wallet, len(rows))
	for i, row := range rows {
		wg.Add(1)
		go func(i int, row string) {
			defer wg.Done()
			var w *wallet.Wallet
			if split != "" {
				items := strings.Split(row, split)
				if len(items) != 4 {
					panic("item number is not 4")
				}
				var err error
				//some check
				w, err = wallet.NewWallet(items[0])
				if err != nil {
					panic(err)
				}
				if hex.EncodeToString(w.GetPrivateKey()) != items[1] {
					panic(model.ErrPrivateKeyErr)
				}
				if hex.EncodeToString(w.GetPublicKey()) != items[2] {
					panic(model.ErrPublicKeyErr)
				}
				if w.GetAddress() != items[3] {
					panic(model.ErrAddressErr)
				}
			} else {
				var err error
				w, err = wallet.NewWallet(row)
				if err != nil {
					panic(err)
				}
			}
			m.Lock()
			wallets[i] = w
			m.Unlock()
		}(i, row)
	}
	wg.Wait()
	return wallets, nil
}

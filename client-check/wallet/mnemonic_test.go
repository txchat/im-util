package wallet

import (
	"encoding/hex"
	walletapi "github.com/txchat/chatcipher"
	"testing"
)

func Test_NewWallet(t *testing.T) {
	mn, err := walletapi.NewMnemonicString(1, 160)
	if err != nil {
		t.Error(err)
		return
	}
	wl, err := NewWallet(mn)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("address: ", wl.GetAddress())
	t.Log("pub: ", hex.EncodeToString(wl.GetPublicKey()))
	t.Log("pri: ", hex.EncodeToString(wl.GetPrivateKey()))
	t.Log("mn: ", wl.GetMnemonic())
	t.Log("signature", wl.GetToken())
}

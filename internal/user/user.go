package user

import (
	"time"

	"github.com/txchat/dtalk/pkg/auth"
	xcrypt "github.com/txchat/dtalk/pkg/crypt"
	secp256k1_ethereum "github.com/txchat/dtalk/pkg/crypt/secp256k1-ethereum"
)

var (
	driver xcrypt.Encrypt
)

func init() {
	var err error
	driver, err = xcrypt.Load(secp256k1_ethereum.Name)
	if err != nil {
		panic(err)
	}
}

var TokenTimeout = time.Minute * 2

type User struct {
	address        string
	userName       string
	priKey, pubKey []byte

	token       string
	tokenExpire time.Time
}

func (u *User) genToken() string {
	authenticator := auth.NewDefaultApiAuthenticatorAsDriver(driver)
	return authenticator.Request("dtalk", u.pubKey, u.priKey)
}

func (u *User) Token() string {
	if time.Now().After(u.tokenExpire.Add(TokenTimeout)) {
		u.tokenExpire = time.Now()
		return u.genToken()
	}
	return u.token
}

func (u *User) SetUsername(username string) {
	u.userName = username
}

func (u *User) GetUsername() string {
	if u.userName == "" {
		return u.address
	}
	return u.userName
}

func (u *User) GetUID() string {
	return u.address
}

func NewUser(address string, priKey, pubKey []byte) *User {
	return &User{
		address: address,
		priKey:  priKey,
		pubKey:  pubKey,
	}
}

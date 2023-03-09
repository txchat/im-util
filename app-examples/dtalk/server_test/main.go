package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/types"
	"github.com/33cn/chain33/wallet/bipwallet"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/txchat/dtalk/api/proto/auth"
	"github.com/txchat/dtalk/api/proto/message"
	"github.com/txchat/im-util/pkg/device"
	"github.com/txchat/im-util/pkg/net"
	"github.com/txchat/im-util/pkg/user"
)

var (
	log           zerolog.Logger
	appId, server string
	scheme        string
)

func init() {
	log = zerolog.New(os.Stdout).With().Timestamp().Logger()
	flag.StringVar(&appId, "appId", "dtalk", "")
	flag.StringVar(&server, "server", "127.0.0.1:3102", "")
	flag.StringVar(&scheme, "scheme", "ws", "tcp/ws")
}

func main() {
	flag.Parse()
	log.Info().Str("appId", appId).Str("server", server).Msg("start")
	var sendAction device.ActionInfo

	a, err := newUser()
	if err != nil {
		log.Error().Err(err).Msg("create user a failed")
		os.Exit(1)
	}

	devA, err := dial(fmt.Sprintf("%v://%v", scheme, server), a, func(c *net.IMConn, action device.ActionInfo) error {
		log.Info().Msg("发送成功，步骤1达成")
		sendAction = action
		return nil
	}, func(c *net.IMConn, action device.ActionInfo) error {
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("create device a failed")
		os.Exit(1)
	}

	b, err := newUser()
	if err != nil {
		log.Error().Err(err).Msg("create user b failed")
		os.Exit(1)
	}

	devB, err := dial(fmt.Sprintf("%v://%v", scheme, server), b, func(c *net.IMConn, action device.ActionInfo) error {
		return nil
	}, func(c *net.IMConn, action device.ActionInfo) error {
		if isFetch(sendAction, action) {
			log.Info().Msg("对端接收成功，步骤2达成")
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("create device b failed")
		os.Exit(1)
	}

	log.Info().Msg(`
	待确认2个步骤：
步骤1：a成功发送消息
步骤2：b接收到a发送的数据
`)

	err = devB.SendTextMsg(message.Channel_Private, a.GetUID(), "hello")
	if err != nil {
		log.Error().Err(err).Msg("user a send proto failed")
		os.Exit(1)
	}

	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			devA.TurnOff()
			devB.TurnOff()
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}

func newUser() (*user.User, error) {
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
	private, public, err := wallet.NewKeyPair(0)
	if err != nil {
		return nil, err
	}
	addr := address.PubKeyToAddr(0, public)
	return user.NewUser(addr, private, public), nil
}

func dial(rawURL string, u *user.User, sendCB device.OnSendHandler, revCB device.OnReceiveHandler) (*device.Device, error) {
	URL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	dev := device.NewDevice(uuid.NewString(), "test_server_a", auth.Device_Android, u)
	err = dev.DialIMServer(appId, *URL, nil)
	if err != nil {
		return nil, err
	}
	dev.SetOnSend(sendCB)
	dev.SetOnReceive(revCB)
	return dev, dev.TurnOn()
}

func isFetch(sendAction, revAction device.ActionInfo) bool {
	if sendAction.Mid == revAction.Mid && sendAction.ChannelType == revAction.ChannelType && sendAction.Target == revAction.Target {
		return true
	}
	return false
}

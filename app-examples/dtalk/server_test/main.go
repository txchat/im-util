package main

import (
	"flag"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	log               zerolog.Logger
	appId             string
	cometAPI, chatAPI string
)

func init() {
	log = zerolog.New(os.Stdout).With().Timestamp().Logger()
	flag.StringVar(&appId, "appId", "dtalk", "")
	flag.StringVar(&cometAPI, "comet", "ws://127.0.0.1:3102", "ws://<IP>:<port> or tcp://<IP>:<port>")
	flag.StringVar(&chatAPI, "chat", "http://127.0.0.1:8888", "http://<IP>:<port> or https://<IP>:<port>")
}

func main() {
	flag.Parse()
	log.Info().Str("appId", appId).Str("cometAPI", cometAPI).Str("chatAPI", chatAPI).Msg("start")
	var sendAction device.ActionInfo

	a, err := newUser()
	if err != nil {
		log.Error().Err(err).Msg("create user a failed")
		os.Exit(1)
	}

	devA, err := dial(chatAPI, cometAPI, a, func(c *net.IMConn, action device.ActionInfo) error {
		if action.Err != nil {
			log.Error().Err(action.Err).Msg("步骤1发送失败")
			return err
		}
		sendAction = action
		log.Info().Str("mid", action.Mid).Msg("发送成功，步骤1达成")
		return nil
	}, nil)
	if err != nil {
		log.Error().Err(err).Msg("create device a failed")
		os.Exit(1)
	}
	log.Info().Str("uid", a.GetUID()).Msg("user a init")

	b, err := newUser()
	if err != nil {
		log.Error().Err(err).Msg("create user b failed")
		os.Exit(1)
	}

	devB, err := dial(chatAPI, cometAPI, b, nil, func(c *net.IMConn, action device.ActionInfo) error {
		if isFetch(sendAction, action) {
			log.Info().Msg("对端接收成功，步骤2达成")
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("create device b failed")
		os.Exit(1)
	}
	log.Info().Str("uid", b.GetUID()).Msg("user b init")

	log.Info().Msg(`
	待确认2个步骤：
步骤1：a成功发送消息
步骤2：b接收到a发送的数据
`)

	err = devA.SendTextMsg(message.Channel_Private, b.GetUID(), "hello")
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

func dial(chatURLStr, cometURLStr string, u *user.User, sendCB device.OnSendHandler, revCB device.OnReceiveHandler) (*device.Device, error) {
	cometURL, err := url.Parse(cometURLStr)
	if err != nil {
		return nil, err
	}
	chatURL, err := url.Parse(chatURLStr)
	if err != nil {
		return nil, err
	}
	dev := device.NewDevice(uuid.NewString(), "test_server_a", auth.Device_Android, u)
	err = dev.DialIMServer(appId, *cometURL, nil)
	if err != nil {
		return nil, err
	}
	dev.DialChatAPI(*chatURL, time.Second*5)
	dev.SetOnSend(sendCB)
	dev.SetOnReceive(revCB)
	return dev, dev.TurnOn()
}

func isFetch(sendAction, revAction device.ActionInfo) bool {
	log.Debug().Interface("sendAction", sendAction).Interface("revAction", revAction).Msg("do isFetch")
	if sendAction.Mid == revAction.Mid && sendAction.ChannelType == revAction.ChannelType && sendAction.Target == revAction.Target {
		return true
	}
	return false
}

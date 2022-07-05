package keep

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	xlog "github.com/txchat/im-util/internel/log"
	"github.com/txchat/im-util/internel/rate"
	"github.com/txchat/im-util/internel/reader"
	"github.com/txchat/im-util/internel/user"
	"github.com/txchat/im-util/protocol/wallet"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var Cmd = &cobra.Command{
	Use:   "keep",
	Short: "",
	Long:  ``,
	Run:   keepRun,
}

var (
	userNum   int
	server    string
	appId     string
	totalTime string

	userStorePath string
	readSplit     string
)

func init() {
	Cmd.Flags().IntVarP(&userNum, "users", "u", 2, "users number")
	Cmd.Flags().StringVarP(&server, "server", "s", "172.16.101.107:3102", "server address")
	Cmd.Flags().StringVarP(&appId, "appId", "a", "dtalk", "")
	Cmd.Flags().StringVarP(&userStorePath, "in", "i", "./users.txt", "users store file path")
	Cmd.Flags().StringVarP(&readSplit, "rs", "", ",", "存储用户信息的字段分隔符[默认：,]")
	Cmd.Flags().StringVarP(&totalTime, "time", "t", "20s", "")
}

func keepRun(cmd *cobra.Command, args []string) {
	start := time.Now()
	//load users
	log := xlog.NewLogger(os.Stdout)
	log.Info().Str("server", server).
		Str("appId", appId).
		Str("totalTime", totalTime).
		Str("userStorePath", userStorePath).
		Int("userNum", userNum).Msg("success config")
	log.Info().Msg("start load user store")

	ttTime, err := rate.ParseDuration(totalTime)
	if err != nil {
		log.Error().Err(err).Msg("ParseRateString error")
		return
	}

	//读取用户信息文件，为了加快生成速度文件存储完整的助记词、私钥、公钥、地址
	metadata, err := reader.LoadMetadata(userStorePath, readSplit)
	if err != nil {
		log.Error().Err(err).Msg("ReadFile error")
		return
	}
	log.Info().Msg(fmt.Sprintf("success load users:%d", len(metadata)))
	if len(metadata) < userNum {
		//TODO 报错
		return
	}

	log.Info().Msg("start create user store")
	var users []*user.User
	m := sync.Mutex{}
	wg := sync.WaitGroup{}
	for _, md := range metadata[:userNum] {
		wg.Add(1)
		go func(md *wallet.Metadata) {
			//defer func() {
			//	if r := recover(); r != nil {
			//		log.Error().Str("address", w.Address).
			//			Interface("priv key", hex.EncodeToString(w.PrivKey)).
			//			Interface("pub key", hex.EncodeToString(w.PubKey)).
			//			Msg("panic create user")
			//		panic(r)
			//	}
			//}()
			u := user.NewUser(md.GetAddress(), md.GetPrivateKey(), md.GetPublicKey())
			m.Lock()
			users = append(users, u)
			m.Unlock()
			wg.Done()
		}(md)
	}
	wg.Wait()
	for _, u := range users {
		err := u.ConnServer(appId, server)
		if err != nil {
			log.Error().Err(err).Msg("Create User error")
			return
		}
	}

	log.Info().Int("users", len(users)).Str("cost", fmt.Sprintln(time.Since(start))).Msg("create users success")

	ctx, closer := context.WithTimeout(context.Background(), ttTime)
	defer closer()
	//block
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		var s os.Signal
		select {
		case s = <-c:
		case <-ctx.Done():
			s = syscall.SIGQUIT
		}
		log.Info().Str("signal", s.String()).Msg("service get a signal")
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			//close
			log.Info().Msg("range send stopped, wait receive follow-up message")
			for _, u := range users {
				u.Close()
			}
			log.Info().Msg("all job down")
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}

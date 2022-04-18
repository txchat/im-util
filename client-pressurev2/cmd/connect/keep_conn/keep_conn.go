package keep_conn

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/txchat/im-util/client-pressurev2/pkg/filehelper"
	"github.com/txchat/im-util/client-pressurev2/pkg/logger"
	"github.com/txchat/im-util/client-pressurev2/pkg/rate"
	"github.com/txchat/im-util/client-pressurev2/pkg/user"
	"github.com/txchat/im-util/client-pressurev2/pkg/wallet"
)

var Cmd = &cobra.Command{
	Use:     "keep",
	Short:   "",
	Long:    "",
	Example: "",
	Run:     do,
}

var (
	userNum   int
	server    string
	appId     string
	totalTime string

	userStorePath string
	sysLogPath    string
)

func init() {
	Cmd.Flags().IntVarP(&userNum, "users", "u", 2, "users number")
	Cmd.Flags().StringVarP(&server, "server", "s", "172.16.101.107:3102", "server address")
	Cmd.Flags().StringVarP(&appId, "appId", "a", "dtalk", "")
	Cmd.Flags().StringVarP(&userStorePath, "in", "i", "./users.txt", "users store file path")
	Cmd.Flags().StringVarP(&sysLogPath, "syslog", "", "./pressure_sys_log.txt", "")
	Cmd.Flags().StringVarP(&totalTime, "time", "t", "20s", "")
}

func do(cmd *cobra.Command, args []string) {
	sysFd, sysCloser, err := filehelper.WriteFile(sysLogPath)
	if err != nil {
		panic(err)
	}
	defer sysCloser()

	start := time.Now()
	//load users
	log := logger.NewSysLog(sysFd)
	log.Info().Str("server", server).
		Str("appId", appId).
		Str("totalTime", totalTime).
		Str("userStorePath", userStorePath).
		Str("sysLogPath", sysLogPath).
		Int("userNum", userNum).Msg("success config")
	log.Info().Msg("start load user store")

	ttTime, err := rate.ParseDuration(totalTime)
	if err != nil {
		log.Error().Err(err).Msg("ParseRateString error")
		return
	}

	fr := filehelper.NewFileReader()
	err = fr.ReadFile(userStorePath)
	if err != nil {
		log.Error().Err(err).Msg("ReadFile error")
		return
	}

	log.Info().Msg("start create user store")
	var users []*user.User
	m := sync.Mutex{}
	wallets := fr.GetUserWallet(userNum)
	log.Info().Msg(fmt.Sprintf("success init user num:%d", len(wallets)))
	wg := sync.WaitGroup{}
	for _, w := range wallets {
		wg.Add(1)
		go func(w *wallet.Wallet) {
			//defer func() {
			//	if r := recover(); r != nil {
			//		log.Error().Str("address", w.Address).
			//			Interface("priv key", hex.EncodeToString(w.PrivKey)).
			//			Interface("pub key", hex.EncodeToString(w.PubKey)).
			//			Msg("panic create user")
			//		panic(r)
			//	}
			//}()
			u := user.NewUser(w.Address, w.PrivKey, w.PubKey)
			m.Lock()
			users = append(users, u)
			m.Unlock()
			wg.Done()
		}(w)
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

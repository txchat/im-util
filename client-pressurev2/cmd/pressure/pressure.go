package pressure

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/txchat/im-util/client-pressurev2/pkg/filehelper"
	"github.com/txchat/im-util/client-pressurev2/pkg/logger"
	"github.com/txchat/im-util/client-pressurev2/pkg/msggenerator"
	"github.com/txchat/im-util/client-pressurev2/pkg/rate"
	"github.com/txchat/im-util/client-pressurev2/pkg/user"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var Cmd = &cobra.Command{
	Use:     "pre",
	Short:   "",
	Long:    "",
	Example: "",
	Run:     do,
}

var (
	userNum   int
	server    string
	appId     string
	rateStr   string
	totalTime string

	userStorePath string
	outputPath    string
	sysLogPath    string
)

func init() {
	Cmd.Flags().IntVarP(&userNum, "users", "u", 2, "users number")
	Cmd.Flags().StringVarP(&server, "server", "s", "172.16.101.107:3102", "server address")
	Cmd.Flags().StringVarP(&appId, "appId", "a", "dtalk", "")
	Cmd.Flags().StringVarP(&outputPath, "out", "o", "./pressure_output.txt", "")
	Cmd.Flags().StringVarP(&userStorePath, "in", "i", "./users.txt", "users store file path")
	Cmd.Flags().StringVarP(&sysLogPath, "syslog", "", "./pressure_sys_log.txt", "")
	Cmd.Flags().StringVarP(&rateStr, "rate", "r", "1/s", "")
	Cmd.Flags().StringVarP(&totalTime, "time", "t", "20s", "")
}

func do(cmd *cobra.Command, args []string) {
	sysFd, sysCloser, err := filehelper.WriteFile(sysLogPath)
	if err != nil {
		panic(err)
	}
	defer sysCloser()
	msgFd, msgCloser, err := filehelper.WriteFile(outputPath)
	if err != nil {
		panic(err)
	}
	defer msgCloser()

	//load users
	log := logger.NewSysLog(sysFd)
	num, tm, err := rate.ParseRateString(rateStr)
	if err != nil {
		log.Error().Err(err).Msg("ParseRateString error")
		return
	}
	ttTime, err := rate.ParseDuration(totalTime)
	if err != nil {
		log.Error().Err(err).Msg("ParseRateString error")
		return
	}

	log.Info().Str("server", server).
		Str("appId", appId).
		Str("rateStr", rateStr).
		Str("totalTime", totalTime).
		Str("userStorePath", userStorePath).
		Str("outputPath", outputPath).
		Str("sysLogPath", sysLogPath).
		Int("userNum", userNum).Msg("success config")
	log.Info().Msg("start load user store")
	fr := filehelper.NewFileReader()
	err = fr.ReadFile(userStorePath)
	if err != nil {
		log.Error().Err(err).Msg("ReadFile error")
		return
	}

	log.Info().Msg("start create user store")
	var users []*user.User
	wallets := fr.GetUserWallet(userNum)
	log.Info().Msg(fmt.Sprintf("success init user num:%d", len(wallets)))
	for _, wallet := range wallets {
		u := user.NewUser(wallet.Address, wallet.PrivKey, wallet.PubKey)
		err := u.ConnServer(appId, server)
		if err != nil {
			log.Error().Err(err).Msg("Create User error")
			continue
		}
		users = append(users, u)
	}
	log.Info().Msg("create users success")

	mg := msggenerator.NewMsgGenerator(users, logger.NewMsgLog(msgFd))

	inv := time.Duration(int(tm) / num)
	log.Info().Msg(fmt.Sprintf("start range send, %s interval pre message", inv.String()))
	for _, u := range users {
		go mg.RangeSend(u, inv, log)
		go mg.HandleAck(u, log)
	}

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
			mg.StopSend()
			log.Info().Msg("range send stopped, wait receive follow-up message")
			time.Sleep(time.Second * 30)
			mg.StopAck()
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

package pressure

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/txchat/im-util/internal/device"
	xlog "github.com/txchat/im-util/internal/log"
	"github.com/txchat/im-util/internal/rate"
	"github.com/txchat/im-util/internal/reader"
	"github.com/txchat/im-util/internal/user"
	"github.com/txchat/im-util/pkg/util"
	"github.com/txchat/im-util/pressure/pkg/msggenerator"
)

var Cmd = &cobra.Command{
	Use:   "pre",
	Short: "提供批量压测功能，将待分析结果输出到指定文件中",
	Long:  `提供批量压测功能，将待分析结果输出到指定文件中。需要配和analyze分析工具，将输出结果分析成可读数据。`,
	RunE:  pressureRunE,
}

var (
	userNum   int
	server    string
	appId     string
	rateStr   string
	totalTime string

	userStorePath string
	readSplit     string
	outputPath    string
)

func init() {
	Cmd.Flags().IntVarP(&userNum, "users", "u", 2, "users number")
	Cmd.Flags().StringVarP(&server, "server", "s", "172.16.101.107:3102", "server address")
	Cmd.Flags().StringVarP(&appId, "appId", "a", "dtalk", "")
	Cmd.Flags().StringVarP(&outputPath, "out", "o", "./pressure_output.txt", "")
	Cmd.Flags().StringVarP(&userStorePath, "in", "i", "./users.txt", "users store file path")
	Cmd.Flags().StringVarP(&readSplit, "rs", "", ",", "存储用户信息的字段分隔符[默认：,]")
	Cmd.Flags().StringVarP(&rateStr, "rate", "r", "1/s", "")
	Cmd.Flags().StringVarP(&totalTime, "time", "t", "20s", "")
}

func pressureRunE(cmd *cobra.Command, args []string) error {
	// 打开文件
	fd, closer, err := util.WriteFile(outputPath)
	if err != nil {
		return err
	}
	defer closer()
	outLog := xlog.NewLogger(fd)
	log := xlog.NewLogger(os.Stdout)

	// rate
	num, tm, err := rate.ParseRateString(rateStr)
	if err != nil {
		return fmt.Errorf("ParseRateString failed: %v", err)
	}
	ttTime, err := rate.ParseDuration(totalTime)
	if err != nil {
		return fmt.Errorf("ParseDuration failed: %v", err)
	}

	log.Info().Str("server", server).
		Str("appId", appId).
		Str("rateStr", rateStr).
		Str("totalTime", totalTime).
		Str("userStorePath", userStorePath).
		Str("outputPath", outputPath).
		Int("userNum", userNum).Msg("success config")
	log.Info().Msg("config")

	//读取用户信息文件，为了加快生成速度文件存储完整的助记词、私钥、公钥、地址
	metadata, err := reader.LoadMetadata(userStorePath, readSplit)
	if err != nil {
		return fmt.Errorf("LoadMetadata failed: %v", err)
	}
	if len(metadata) < userNum {
		log.Error().Err(err).Int("len(metadata)", len(metadata)).Int("userNum", userNum).Msg("users lacking")
		return fmt.Errorf("length of metadata less than userNum")
	}
	log.Info().Int("len(metadata)", len(metadata)).Int("userNum", userNum).Msg("LoadMetadata")

	var users []*user.User
	for _, md := range metadata[:userNum] {
		u := user.NewUser(md.GetAddress(), md.GetPrivateKey(), md.GetPublicKey())
		users = append(users, u)
	}

	var devices []*device.Device
	for _, u := range users {
		d := device.NewDevice("", "", 0, outLog, u)
		err = d.DialIMServer(appId, server, nil)
		if err != nil {
			log.Error().Err(err).Msg("DialIMServer failed")
			continue
		}
		err = d.TurnOn()
		if err != nil {
			log.Error().Err(err).Msg("Device TurnOn failed")
			continue
		}
		devices = append(devices, d)
	}
	log.Info().Msg("all init success!")

	mg := msggenerator.NewMsgGenerator(users)

	inv := time.Duration(int(tm) / num)
	log.Info().Msg(fmt.Sprintf("start range send, %s interval pre message", inv.String()))
	for _, d := range devices {
		go mg.RangeSend(d, inv)
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
			for _, d := range devices {
				d.TurnOff()
			}
			log.Info().Msg("all job down")
			return nil
		case syscall.SIGHUP:
			// TODO reload
		default:
			return nil
		}
	}
}

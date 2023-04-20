package pressure

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/txchat/dtalk/api/proto/auth"
	xlog "github.com/txchat/im-util/internal/log"
	"github.com/txchat/im-util/internal/rate"
	"github.com/txchat/im-util/pkg/device"
	"github.com/txchat/im-util/pkg/net"
	"github.com/txchat/im-util/pkg/user"
	"github.com/txchat/im-util/pkg/util"
	"github.com/txchat/im-util/pkg/wallet"
	"github.com/txchat/im-util/pressure/internal/msggenerator"
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
	appID     string
	rateStr   string
	totalTime string

	userStorePath string
	readSplit     string
	outputPath    string
)

func init() {
	Cmd.Flags().IntVarP(&userNum, "users", "u", 2, "users number")
	Cmd.Flags().StringVarP(&server, "server", "s", "ws://172.16.101.107:3102", "server address")
	Cmd.Flags().StringVarP(&appID, "appId", "a", "dtalk", "")
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
		Str("appId", appID).
		Str("rateStr", rateStr).
		Str("totalTime", totalTime).
		Str("userStorePath", userStorePath).
		Str("outputPath", outputPath).
		Int("userNum", userNum).Msg("success config")
	log.Info().Msg("config")

	//读取用户信息文件，为了加快生成速度文件存储完整的助记词、私钥、公钥、地址
	metadata, err := wallet.LoadMetadata(userStorePath, readSplit)
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

	URL, err := url.Parse(server)
	if err != nil {
		return fmt.Errorf("url parse failed: %v", err)
	}

	lp := NewLogPrinter(outLog)
	var devices []*device.Device
	for _, u := range users {
		d := device.NewDevice(uuid.NewString(), "pressure-test", auth.Device_Android, u)
		err = d.DialIMServer(appID, *URL, nil)
		if err != nil {
			log.Error().Err(err).Msg("DialIMServer failed")
			continue
		}
		d.SetOnSend(lp.onSendLogs)
		d.SetOnReceive(lp.onReceiveLogs)
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

type LogPrinter struct {
	log zerolog.Logger
}

func NewLogPrinter(log zerolog.Logger) *LogPrinter {
	return &LogPrinter{log: log}
}

func (lp *LogPrinter) onSendLogs(c *net.IMConn, action device.ActionInfo) error {
	// 发出时间点的日志
	lp.log.Info().Str("action", action.Action).
		Str("user_id", action.UID).
		Str("conn_id", action.ConnID).
		Str("uuid", action.UUID).
		Str("from", action.From).
		Str("target", action.Target).
		Str("channel_type", action.ChannelType.String()).
		Int32("seq", action.Seq).
		Int32("ack", action.Ack).
		Str("mid", action.Mid).
		Msg("")
	return nil
}

func (lp *LogPrinter) onReceiveLogs(c *net.IMConn, action device.ActionInfo) error {
	// 接收时间点的日志
	lp.log.Info().Str("action", action.Action).
		Str("user_id", action.UID).
		Str("conn_id", action.ConnID).
		Str("uuid", action.UUID).
		Str("from", action.From).
		Str("target", action.Target).
		Str("channel_type", action.ChannelType.String()).
		Int32("seq", action.Seq).
		Int32("ack", action.Ack).
		Str("mid", action.Mid).
		Msg("")
	return nil
}

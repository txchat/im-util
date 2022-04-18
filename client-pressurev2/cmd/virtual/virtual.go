package virtual

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	xgrpc "github.com/txchat/dtalk/pkg/net/grpc"
	xtime "github.com/txchat/dtalk/pkg/time"
	"github.com/txchat/im-util/client-pressurev2/pkg/device"
	"github.com/txchat/im-util/client-pressurev2/pkg/filehelper"
	xproto "github.com/txchat/imparse/proto"
)

var Cmd = &cobra.Command{
	Use:     "device",
	Short:   "",
	Long:    "",
	Example: "",
	Run:     do,
}

var (
	userNum    int
	server     string
	appId      string
	uuid       string
	deviceName string
	deviceType int32
	port       int32

	userStorePath string
)

func init() {
	Cmd.Flags().IntVarP(&userNum, "users", "u", 1, "users number")
	Cmd.Flags().StringVarP(&server, "server", "s", "172.16.101.107:3102", "server address")
	Cmd.Flags().StringVarP(&appId, "appId", "a", "dtalk", "")
	Cmd.Flags().StringVarP(&uuid, "uuid", "", "3ade6a21-a0d7-48ce-94a2-2f3567adc468", "")
	Cmd.Flags().StringVarP(&deviceName, "dname", "", "虚拟驱动", "")
	Cmd.Flags().Int32VarP(&deviceType, "dtype", "", 0, "")
	Cmd.Flags().Int32VarP(&port, "port", "p", 30001, "")
	Cmd.Flags().StringVarP(&userStorePath, "in", "i", "./users.txt", "users store file path")
}

func do(cmd *cobra.Command, args []string) {
	//load users
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	log.Info().Str("server", server).
		Str("appId", appId).
		Str("uuid", uuid).
		Str("deviceName", deviceName).
		Int32("deviceType", deviceType).
		Str("userStorePath", userStorePath).
		Int("userNum", userNum).Msg("success config")
	log.Info().Msg("start load user store")

	fr := filehelper.NewFileReader()
	err := fr.ReadFile(userStorePath)
	if err != nil {
		log.Error().Err(err).Msg("ReadFile error")
		return
	}

	log.Info().Msg("start create wallet")
	wallets := fr.GetUserWallet(userNum)
	log.Info().Msg(fmt.Sprintf("success init user num:%d", len(wallets)))

	if len(wallets) < 1 {
		log.Error().Msg(fmt.Sprintf("wallets numer is %d", len(wallets)))
		return
	}

	wallet := wallets[0]

	d := device.NewNetScreenDriver(&device.UserDeviceOpt{
		AppId:      appId,
		Server:     server,
		Address:    wallet.Address,
		PriKey:     wallet.PrivKey,
		PubKey:     wallet.PubKey,
		Uuid:       uuid,
		DeviceName: deviceName,
		DeviceType: xproto.Device(deviceType),
	}, &xgrpc.ServerConfig{
		Network:                           "tcp",
		Addr:                              fmt.Sprintf(":%d", port),
		Timeout:                           xtime.Duration(time.Second),
		KeepAliveMaxConnectionIdle:        xtime.Duration(time.Second * 60),
		KeepAliveMaxConnectionAge:         xtime.Duration(time.Hour * 2),
		KeepAliveMaxMaxConnectionAgeGrace: xtime.Duration(time.Second * 20),
		KeepAliveTime:                     xtime.Duration(time.Second * 60),
		KeepAliveTimeout:                  xtime.Duration(time.Second * 20),
	})
	d.StartUp()
	log.Info().Msg("success on serve")

	//block
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info().Str("signal", s.String()).Msg("service get a signal")
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			//close
			ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
			d.Shutdown(ctx)
			log.Info().Msg("all job down")
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}

package device

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	xgrpc "github.com/txchat/dtalk/pkg/net/grpc"
	xtime "github.com/txchat/dtalk/pkg/time"
	"github.com/txchat/im-util/internel/device"
	"github.com/txchat/im-util/internel/user"
	"github.com/txchat/im-util/pressure/internel/reader"
	deviceGRPC "github.com/txchat/im-util/virtual/grpc/device"
	"github.com/txchat/imparse/proto"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var Cmd = &cobra.Command{
	Use:   "device",
	Short: "客户端设备，与聊天服务建立通讯",
	Long:  ``,
	Run:   deviceRun,
}

var (
	server     string
	appId      string
	uuid       string
	deviceName string
	deviceType int32
	port       int32

	userStorePath string
	readSplit     string
)

func init() {
	Cmd.Flags().StringVarP(&server, "server", "s", "172.16.101.107:3102", "服务端地址")
	Cmd.Flags().StringVarP(&appId, "appId", "a", "dtalk", "应用ID")
	Cmd.Flags().StringVarP(&uuid, "uuid", "", "3ade6a21-a0d7-48ce-94a2-2f3567adc468", "设备唯一识别号")
	Cmd.Flags().StringVarP(&deviceName, "dname", "", "虚拟驱动", "设备名称")
	Cmd.Flags().Int32VarP(&deviceType, "dtype", "", 0, "设备类型:[0]Android")
	Cmd.Flags().Int32VarP(&port, "port", "p", 30001, "")
	Cmd.Flags().StringVarP(&userStorePath, "in", "i", "./users.txt", "users store file path")
	Cmd.Flags().StringVarP(&readSplit, "rs", "", ",", "存储用户信息的字段分隔符[默认：,]")
}

func deviceRun(cmd *cobra.Command, args []string) {
	//load users
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	log.Info().Str("server", server).
		Str("appId", appId).
		Str("uuid", uuid).
		Str("deviceName", deviceName).
		Int32("deviceType", deviceType).
		Str("userStorePath", userStorePath).Msg("success config")
	log.Info().Msg("start load user store")

	//读取用户信息文件，为了加快生成速度文件存储完整的助记词、私钥、公钥、地址
	metadata, err := reader.LoadMetadata(userStorePath, readSplit)
	if err != nil {
		log.Error().Err(err).Msg("ReadFile error")
		return
	}
	log.Info().Msg(fmt.Sprintf("success load users:%d", len(metadata)))
	if len(metadata) < 1 {
		log.Error().Msg(fmt.Sprintf("load users numer is %d", len(metadata)))
		return
	}
	md := metadata[0]

	u := user.NewUser(md.GetAddress(), md.GetPrivateKey(), md.GetPublicKey())
	d := device.NewDevice(uuid, deviceName, proto.Device(deviceType), log, u)
	err = d.DialIMServer(appId, server, nil)
	if err != nil {
		log.Error().Err(err).Msg("DialIMServer failed")
		return
	}
	err = d.TurnOn()
	if err != nil {
		log.Error().Err(err).Msg("Device TurnOn failed")
		return
	}
	srvGrpc := deviceGRPC.NewServer(&xgrpc.ServerConfig{
		Network:                           "tcp",
		Addr:                              fmt.Sprintf(":%d", port),
		Timeout:                           xtime.Duration(time.Second),
		KeepAliveMaxConnectionIdle:        xtime.Duration(time.Second * 60),
		KeepAliveMaxConnectionAge:         xtime.Duration(time.Hour * 2),
		KeepAliveMaxMaxConnectionAgeGrace: xtime.Duration(time.Second * 20),
		KeepAliveTime:                     xtime.Duration(time.Second * 60),
		KeepAliveTimeout:                  xtime.Duration(time.Second * 20),
	}, d)
	log.Info().Msg("success on serve")

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
			d.TurnOff()
			srvGrpc.Shutdown(ctx)
			log.Info().Msg("all job down")
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}

/*
Copyright © 2022 oofpgDLD <oofpgdld@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package connect

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/txchat/im-util/internal/device"
	xlog "github.com/txchat/im-util/internal/log"
	"github.com/txchat/im-util/internal/rate"
	"github.com/txchat/im-util/internal/reader"
	"github.com/txchat/im-util/internal/user"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// keepCmd represents the keep command
var keepCmd = &cobra.Command{
	Use:   "keep",
	Short: "启动指定数量的客户端连接",
	Long:  `启动指定数量的客户端连接，并保持心跳，不发送消息`,
	RunE:  keepRunE,
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
	Cmd.AddCommand(keepCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// keepCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// keepCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	keepCmd.Flags().IntVarP(&userNum, "users", "u", 2, "users number")
	keepCmd.Flags().StringVarP(&server, "server", "s", "172.16.101.107:3102", "server address")
	keepCmd.Flags().StringVarP(&appId, "appId", "a", "dtalk", "")
	keepCmd.Flags().StringVarP(&userStorePath, "in", "i", "./users.txt", "users store file path")
	keepCmd.Flags().StringVarP(&readSplit, "rs", "", ",", "存储用户信息的字段分隔符[默认：,]")
	keepCmd.Flags().StringVarP(&totalTime, "time", "t", "720h", "")
}

func keepRunE(cmd *cobra.Command, args []string) error {
	start := time.Now()
	//load users
	log := xlog.NewLogger(os.Stdout)
	log.Info().Str("server", server).
		Str("appId", appId).
		Str("totalTime", totalTime).
		Str("userStorePath", userStorePath).
		Int("userNum", userNum).Msg("success config")
	logNil := zerolog.New(nil)

	ttTime, err := rate.ParseDuration(totalTime)
	if err != nil {
		return fmt.Errorf("ParseDuration failed: %v", err)
	}

	//读取用户信息文件，为了加快生成速度文件存储完整的助记词、私钥、公钥、地址
	metadata, err := reader.LoadMetadata(userStorePath, readSplit)
	if err != nil {
		return fmt.Errorf("LoadMetadata failed: %v", err)
	}
	log.Info().Msg(fmt.Sprintf("success load users:%d", len(metadata)))
	if len(metadata) < userNum {
		log.Error().Err(err).Int("len(metadata)", len(metadata)).Int("userNum", userNum).Msg("users lacking")
		return fmt.Errorf("length of metadata less than userNum")
	}
	log.Info().Str("cost", fmt.Sprint(time.Since(start))).Int("len(metadata)", len(metadata)).Int("userNum", userNum).Msg("LoadMetadata")
	start = time.Now()

	//var users []*user.User
	var devices []*device.Device
	for _, md := range metadata[:userNum] {
		u := user.NewUser(md.GetAddress(), md.GetPrivateKey(), md.GetPublicKey())
		//users = append(users, u)
		d := device.NewDevice("", "", 0, logNil, u)
		devices = append(devices, d)
	}
	log.Info().Str("cost", fmt.Sprint(time.Since(start))).Msg("all init success!")
	start = time.Now()

	wg := sync.WaitGroup{}
	for _, dev := range devices {
		wg.Add(1)
		go func(d *device.Device) {
			defer wg.Done()
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
		}(dev)
	}
	wg.Wait()
	log.Info().Int("users", len(devices)).Str("cost", fmt.Sprint(time.Since(start))).Msg("all device connected")

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

package keyboard

import (
	"bufio"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	pb "github.com/txchat/im-util/virtual/grpc/device/api"
)

var Cmd = &cobra.Command{
	Use:     "keyboard",
	Short:   "",
	Long:    "",
	Example: "",
	Run:     do,
}

var (
	server string
)

func init() {
	Cmd.Flags().StringVarP(&server, "server", "s", "localhost:30001", "server address")
}

func do(cmd *cobra.Command, args []string) {
	//load users
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	log.Info().Str("server", server).Msg("success config")

	log.Info().Msg("key board start")
	client := pb.New(server)
	stream, err := client.Input(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("ListFile")
		return
	}

	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			res, _ := reader.ReadString('\n')

			err = stream.Send(&pb.InputReq{
				Text: res,
			})
			if err != nil {
				log.Error().Err(err).Msg("stream.Send")
				return
			}
		}
	}()

	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info().Str("signal", s.String()).Msg("service get a signal")
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			//close
			//关闭流并获取返回的消息
			_, err := stream.CloseAndRecv()
			if err != nil {
				log.Error().Err(err).Msg("stream.CloseAndRecv")
				return
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

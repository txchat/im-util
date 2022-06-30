package screen

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	pb "github.com/txchat/im-util/pressure/pkg/device/api"
	"io"
	"os"
	"os/signal"
	"syscall"
)

var Cmd = &cobra.Command{
	Use:     "screen",
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

	log.Info().Msg("screen start")
	client := pb.New(server)
	stream, err := client.Output(context.Background(), &pb.OutputReq{})
	if err != nil {
		log.Error().Err(err).Msg("ListFile")
		return
	}

	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Error().Err(err).Msg("Recv")
				return
			}
			fmt.Println(msg.GetMsg())
		}
		log.Info().Msg("device power off")
	}()

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
			err = stream.CloseSend()
			if err != nil {
				log.Error().Err(err).Msg("CloseSend")
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

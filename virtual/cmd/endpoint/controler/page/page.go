package page

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	pb "github.com/txchat/im-util/pressure/pkg/device/api"
)

var Cmd = &cobra.Command{
	Use:     "page",
	Short:   "",
	Long:    "",
	Example: "",
	Run:     do,
}

var (
	server string

	toGroup bool
	toUser  bool
	target  string
)

func init() {
	Cmd.Flags().StringVarP(&server, "server", "s", "localhost:30001", "server address")
	Cmd.Flags().StringVarP(&target, "target", "t", "", "message target")
	Cmd.Flags().BoolVarP(&toGroup, "group", "g", false, "")
	Cmd.Flags().BoolVarP(&toUser, "user", "u", false, "")
}

func do(cmd *cobra.Command, args []string) {
	//load users
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	log.Info().Str("server", server).Msg("success config")

	chType := int32(0)
	if toUser {
		chType = 0
	} else if toGroup {
		chType = 1
	} else {
		log.Error().Msg("channel type error")
		return
	}

	client := pb.New(server)
	_, err := client.ChangeCurrentPage(context.Background(), &pb.ChangeCurrentPageReq{
		ChannelType: chType,
		Target:      target,
	})
	if err != nil {
		log.Error().Err(err).Msg("ChangeCurrentPage")
		return
	}

	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info().Str("signal", s.String()).Msg("service get a signal")
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			//close
			log.Info().Msg("all job down")
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}

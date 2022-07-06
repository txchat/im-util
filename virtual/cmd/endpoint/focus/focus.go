package focus

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	pb "github.com/txchat/im-util/virtual/grpc/device/api"
)

var Cmd = &cobra.Command{
	Use:     "focus",
	Short:   "",
	Long:    "",
	Example: "",
	Run:     focusRun,
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

	Cmd.MarkFlagsMutuallyExclusive("group", "user")
	Cmd.MarkFlagRequired("target")
}

func focusRun(cmd *cobra.Command, args []string) {
	//load users
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	log.Info().Str("server", server).Msg("")

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
	_, err := client.Focus(context.Background(), &pb.FocusReq{
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

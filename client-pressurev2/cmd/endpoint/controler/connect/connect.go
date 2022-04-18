package connect

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	pb "github.com/txchat/im-util/client-pressurev2/pkg/device/api"
	"os"
)

var Cmd = &cobra.Command{
	Use:     "conn",
	Short:   "",
	Long:    "",
	Example: "",
	Run:     do,
}

var (
	server string

	connType string
)

func init() {
	Cmd.Flags().StringVarP(&server, "server", "s", "localhost:30001", "server address")
	Cmd.Flags().StringVarP(&connType, "conn", "c", "connect", "type of [connect,reconnect,disconnect]")
}

func do(cmd *cobra.Command, args []string) {
	//load users
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	log.Info().Str("server", server).Msg("success config")

	client := pb.New(server)
	switch connType {
	case "connect":
		_, err := client.Connect(context.Background(), &pb.ConnectReq{})
		if err != nil {
			log.Error().Err(err).Msg("Connect")
			return
		}
		log.Info().Msg("success connect to server")
	case "reconnect":
		_, err := client.ReConnect(context.Background(), &pb.ReConnectReq{})
		if err != nil {
			log.Error().Err(err).Msg("ReConnect")
			return
		}
		log.Info().Msg("success reconnect to server")
	case "disconnect":
		_, err := client.DisConnect(context.Background(), &pb.DisConnectReq{})
		if err != nil {
			log.Error().Err(err).Msg("DisConnect")
			return
		}
		log.Info().Msg("success disconnected")
	}
}

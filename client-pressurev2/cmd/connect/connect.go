package connect

import (
	"github.com/spf13/cobra"
	"github.com/txchat/im-util/client-pressurev2/cmd/connect/keep_conn"
)

var Cmd = &cobra.Command{
	Use:     "conn",
	Short:   "",
	Long:    "",
	Example: "",
	Run:     nil,
}

func init() {
	Cmd.AddCommand(keep_conn.Cmd)
}

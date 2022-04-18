package frame

import (
	"github.com/spf13/cobra"
	"github.com/txchat/im-util/client-check/cmd/frame/auth"
)

var Cmd = &cobra.Command{
	Use:     "frame",
	Short:   "frame",
	Long:    "",
	Example: "frame [command]",
	Run:     do,
}

func init() {
	Cmd.AddCommand(auth.AuthCmd)
}

func do(cmd *cobra.Command, args []string) {
}

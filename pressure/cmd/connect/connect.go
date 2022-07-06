package connect

import (
	"github.com/spf13/cobra"
	"github.com/txchat/im-util/pressure/cmd/connect/keep"
)

var Cmd = &cobra.Command{
	Use:     "conn",
	Short:   "",
	Long:    "",
	Example: "",
	Run:     nil,
}

func init() {
	Cmd.AddCommand(keep.Cmd)
}

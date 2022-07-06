package endpoint

import (
	"github.com/spf13/cobra"
	"github.com/txchat/im-util/virtual/cmd/endpoint/focus"
	"github.com/txchat/im-util/virtual/cmd/endpoint/keyboard"
	"github.com/txchat/im-util/virtual/cmd/endpoint/screen"
)

var Cmd = &cobra.Command{
	Use:   "ep",
	Short: "",
	Long:  ``,
}

func init() {
	Cmd.AddCommand(screen.Cmd)
	Cmd.AddCommand(keyboard.Cmd)
	Cmd.AddCommand(focus.Cmd)
}

package endpoint

import (
	"github.com/spf13/cobra"
	"github.com/txchat/im-util/client-pressurev2/cmd/endpoint/controler"
	"github.com/txchat/im-util/client-pressurev2/cmd/endpoint/keyboard"
	"github.com/txchat/im-util/client-pressurev2/cmd/endpoint/screen"
)

var Cmd = &cobra.Command{
	Use:     "ep",
	Short:   "",
	Long:    "",
	Example: "",
	Run:     nil,
}

func init() {
	Cmd.AddCommand(screen.Cmd)
	Cmd.AddCommand(keyboard.Cmd)
	Cmd.AddCommand(controler.Cmd)
}
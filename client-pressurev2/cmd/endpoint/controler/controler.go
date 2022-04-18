package controler

import (
	"github.com/spf13/cobra"
	"github.com/txchat/im-util/client-pressurev2/cmd/endpoint/controler/connect"
	"github.com/txchat/im-util/client-pressurev2/cmd/endpoint/controler/page"
)

var Cmd = &cobra.Command{
	Use:     "ctl",
	Short:   "",
	Long:    "",
	Example: "",
	Run:     nil,
}

func init() {
	Cmd.AddCommand(page.Cmd)
	Cmd.AddCommand(connect.Cmd)
}

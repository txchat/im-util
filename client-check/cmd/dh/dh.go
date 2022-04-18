package dh

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "dh",
	Short:   "dh about",
	Long:    "",
	Example: "dh enc",
	Run:     do,
}

func init() {
	Cmd.AddCommand(Encrypt)
	Cmd.AddCommand(SessionKey)
}

func do(cmd *cobra.Command, args []string) {
}

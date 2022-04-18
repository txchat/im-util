package auth

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "auth",
	Short:   "auth",
	Long:    "",
	Example: "auth [command]",
	Run:     do,
}

func init() {
	Cmd.AddCommand(CheckCmd)
	Cmd.AddCommand(GenCmd)
}

func do(cmd *cobra.Command, args []string) {
}

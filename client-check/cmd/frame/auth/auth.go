package auth

import (
	"github.com/spf13/cobra"
)

var AuthCmd = &cobra.Command{
	Use:     "auth",
	Short:   "",
	Long:    "",
	Example: "auth [command]",
	Run:     doAuthCmd,
}

func init() {
	AuthCmd.AddCommand(GenCmd)
	AuthCmd.AddCommand(CheckCmd)
}

func doAuthCmd(cmd *cobra.Command, args []string) {
}

package mnemonic

import (
	"github.com/spf13/cobra"
	"github.com/txchat/im-util/client-check/cmd/mnemonic/seed"
)

var Cmd = &cobra.Command{
	Use:     "mne",
	Short:   "mnemonic",
	Long:    "",
	Example: "mne [command]",
	Run:     do,
}

func init() {
	Cmd.AddCommand(seed.SeedCmd)
}

func do(cmd *cobra.Command, args []string) {
}

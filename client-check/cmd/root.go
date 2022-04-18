package cmd

import (
	"fmt"
	"github.com/txchat/im-util/client-check/cmd/address"
	"github.com/txchat/im-util/client-check/cmd/auth"
	"github.com/txchat/im-util/client-check/cmd/dh"
	"github.com/txchat/im-util/client-check/cmd/encrypt"
	"github.com/txchat/im-util/client-check/cmd/frame"
	"github.com/txchat/im-util/client-check/cmd/mnemonic"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "imall",
	Short:   "Im platform useful tools",
	Example: "  tools auth -d <data>\n",
}

func init() {
	rootCmd.AddCommand(address.Cmd)
	rootCmd.AddCommand(auth.Cmd)
	rootCmd.AddCommand(dh.Cmd)
	rootCmd.AddCommand(encrypt.Cmd)
	rootCmd.AddCommand(frame.Cmd)
	rootCmd.AddCommand(mnemonic.Cmd)
}

// Execute executes the root command and its subcommands.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

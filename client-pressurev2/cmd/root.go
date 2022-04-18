package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/txchat/im-util/client-pressurev2/cmd/analyze"
	"github.com/txchat/im-util/client-pressurev2/cmd/connect"
	"github.com/txchat/im-util/client-pressurev2/cmd/endpoint"
	"github.com/txchat/im-util/client-pressurev2/cmd/pressure"
	"github.com/txchat/im-util/client-pressurev2/cmd/virtual"
)

var rootCmd = &cobra.Command{
	Use:     "pressurev2",
	Short:   "",
	Example: "pressurev2 auth -d <data>\n",
}

func init() {
	rootCmd.AddCommand(pressure.Cmd)
	rootCmd.AddCommand(analyze.Cmd)
	rootCmd.AddCommand(endpoint.Cmd)
	rootCmd.AddCommand(virtual.Cmd)
	rootCmd.AddCommand(connect.Cmd)
}

// Execute executes the root command and its subcommands.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

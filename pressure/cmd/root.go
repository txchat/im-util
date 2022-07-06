package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/txchat/im-util/pressure/cmd/analyze"
	"github.com/txchat/im-util/pressure/cmd/pressure"
)

var rootCmd = &cobra.Command{
	Use:     "pressure",
	Short:   "",
	Example: "pressure auth -d <data>\n",
}

func init() {
	rootCmd.AddCommand(pressure.Cmd)
	rootCmd.AddCommand(analyze.Cmd)
}

// Execute executes the root command and its subcommands.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

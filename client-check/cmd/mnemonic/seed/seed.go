package seed

import (
	"github.com/spf13/cobra"
)

var SeedCmd = &cobra.Command{
	Use:     "seed",
	Short:   "seed",
	Long:    "",
	Example: "seed [command]",
	Run:     do,
}

func init() {
	SeedCmd.AddCommand(seedEnc)
	SeedCmd.AddCommand(seedDec)
}

func do(cmd *cobra.Command, args []string) {

}

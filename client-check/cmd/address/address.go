package address

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "addr",
	Short: "BTY address about",
	Long:  "address feature: 1.generate",
	Run:   do,
}

func init() {
	Cmd.AddCommand(GenCmd)
}

func do(cmd *cobra.Command, args []string) {
}

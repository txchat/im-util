package encrypt

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "encrypt",
	Short:   "",
	Long:    "sign and verify",
	Example: "encrypt sign",
	Run:     do,
}

func init() {
	Cmd.AddCommand(SignCmd)
	Cmd.AddCommand(VerifyCmd)
}

func do(cmd *cobra.Command, args []string) {
}

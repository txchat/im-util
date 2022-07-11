package wallet

import (
	"testing"

	"github.com/spf13/cobra"
)

func Test_genWallet(t *testing.T) {
	c := cobra.Command{
		Use: "test",
	}
	c.AddCommand(genCmd)
	output, err := executeCommand(&c, "gen", "-d", "-n=1")
	if err != nil {
		t.Error(err)
	}
	t.Log(output)
}

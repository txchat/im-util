package wallet

import (
	"github.com/spf13/cobra"
	"testing"
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

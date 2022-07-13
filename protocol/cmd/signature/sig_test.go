package signature

import (
	"testing"

	"github.com/spf13/cobra"
)

func Test_sigRunE(t *testing.T) {
	c := cobra.Command{
		Use: "test",
	}
	c.AddCommand(sigCmd)
	output, err := executeCommand(&c, "sig",
		"-p", "d8155deca6da46c5015086b9f10adb5c8e992a0c5f62fef50098d7ecabbbb7a9",
		"-m", "hello")
	if err != nil {
		t.Error(err)
	}
	t.Log(output)
}

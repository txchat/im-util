package token

import (
	"testing"

	"github.com/spf13/cobra"
)

func Test_genRunE(t *testing.T) {
	c := cobra.Command{
		Use: "test",
	}
	c.AddCommand(genCmd)
	output, err := executeCommand(&c, "gen",
		"-p", "d8155deca6da46c5015086b9f10adb5c8e992a0c5f62fef50098d7ecabbbb7a9",
		"-P", "031cb7282d22b4a5910bae74172ad2e548bbaf24f076bc247112c1077bb5c40bce",
		"-t", "1656302880000")
	if err != nil {
		t.Error(err)
	}
	t.Log(output)
}

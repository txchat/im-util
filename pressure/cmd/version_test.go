package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func Test_VersionTmpl(t *testing.T) {
	c := cobra.Command{
		Use:     "test",
		Version: version(),
	}
	output, err := executeCommand(&c, "--version")
	if err != nil {
		t.Error(err)
	}
	t.Log(output)
}

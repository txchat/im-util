package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func Test_docRunEMan(t *testing.T) {
	c := cobra.Command{
		Use: "test",
	}
	c.AddCommand(docCmd)
	output, err := executeCommand(&c, "doc", "-d", "../doc/", "-f", "man")
	if err != nil {
		t.Error(err)
	}
	t.Log(output)
}

func Test_docRunEMarkdown(t *testing.T) {
	c := cobra.Command{
		Use: "test",
	}
	c.AddCommand(docCmd)
	output, err := executeCommand(&c, "doc", "-d", "../doc/", "-f", "markdown")
	if err != nil {
		t.Error(err)
	}
	t.Log(output)
}

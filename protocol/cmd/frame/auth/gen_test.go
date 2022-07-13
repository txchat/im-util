package auth

import (
	"testing"

	"github.com/spf13/cobra"
)

func Test_genRunE(t *testing.T) {
	tk := "cWPt9gmAcC93izSmfUOwb7nAXP5y/K4vghzKssyig3tNv/o3pcGRHjjmeB6t0ay3BpGXAVHfhif2NU8wIKB3xgE=#1656303657343*dtalk#031cb7282d22b4a5910bae74172ad2e548bbaf24f076bc247112c1077bb5c40bce"
	c := cobra.Command{
		Use: "test",
	}
	c.AddCommand(genCmd)
	output, err := executeCommand(&c, "gen", "-t", tk)
	if err != nil {
		t.Error(err)
	}
	t.Log(output)
}

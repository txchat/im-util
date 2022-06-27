package signature

import (
	"github.com/spf13/cobra"
	"testing"
)

func Test_verifyRunE(t *testing.T) {
	c := cobra.Command{
		Use: "test",
	}
	c.AddCommand(verifyCmd)
	output, err := executeCommand(&c, "verify",
		"-P", "031cb7282d22b4a5910bae74172ad2e548bbaf24f076bc247112c1077bb5c40bce",
		"-m", "hello",
		"-s", "5e98486b7b7edb721e2484aed9e6a6f69637dd4933e87f0ac28d84edfb02623c1f8bb7cb5cff8a07d4af19123243b2e94de9a954555b566659f0c37300aeb40b01")
	if err != nil {
		t.Error(err)
	}
	t.Log(output)
}

package auth

import (
	"github.com/spf13/cobra"
	"testing"
)

func Test_checkRunE(t *testing.T) {
	data := "AAAAzQAUAAAAAAABAAAAAAAAAAAKBWR0YWxrEq8BY1dQdDlnbUFjQzkzaXpTbWZVT3diN25BWFA1eS9LNHZnaHpLc3N5aWczdE52L28zcGNHUkhqam1lQjZ0MGF5M0JwR1hBVkhmaGlmMk5VOHdJS0IzeGdFPSMxNjU2MzAzNjU3MzQzKmR0YWxrIzAzMWNiNzI4MmQyMmI0YTU5MTBiYWU3NDE3MmFkMmU1NDhiYmFmMjRmMDc2YmMyNDcxMTJjMTA3N2JiNWM0MGJjZQ=="
	c := cobra.Command{
		Use: "test",
	}
	c.AddCommand(checkCmd)
	output, err := executeCommand(&c, "check", "-d", data)
	if err != nil {
		t.Error(err)
	}
	t.Log(output)
}

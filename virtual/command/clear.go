package command

import (
	"io"
	"os/exec"
)

func Clear(w io.Writer) error {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = w
	return cmd.Run()
}

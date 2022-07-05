package tui

import (
	"testing"
)

func Test_newHeaderRegion(t *testing.T) {
	h := newHeaderRegion()
	t.Log(h.View())
}

package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/txchat/im-util/virtual/tui"
)

func main() {
	p := tea.NewProgram(tui.NewMainPage())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

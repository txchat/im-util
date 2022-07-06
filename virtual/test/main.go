package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/txchat/im-util/virtual/tui"
	"log"
)

func main() {
	p := tea.NewProgram(tui.NewMainPage())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

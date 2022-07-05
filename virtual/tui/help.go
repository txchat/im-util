package tui

//
//import (
//	"github.com/charmbracelet/bubbles/help"
//	"github.com/charmbracelet/bubbles/key"
//	tea "github.com/charmbracelet/bubbletea"
//	"github.com/charmbracelet/lipgloss"
//	"strings"
//)
//
//type helpRegion struct {
//	help help.Model
//}
//
//func newHelpRegion() helpRegion {
//	return helpRegion{
//		keys:       keys,
//		help:       help.New(),
//		inputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
//	}
//}
//
//func (m helpRegion) Init() tea.Cmd {
//	return nil
//}
//
//func (m helpRegion) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//	switch msg := msg.(type) {
//	case tea.WindowSizeMsg:
//		// If we set a width on the help menu it can it can gracefully truncate
//		// its view as needed.
//		m.help.Width = msg.Width
//
//	case tea.KeyMsg:
//		switch {
//		case key.Matches(msg, m.keys.Up):
//			m.lastKey = "↑"
//		case key.Matches(msg, m.keys.Down):
//			m.lastKey = "↓"
//		case key.Matches(msg, m.keys.Left):
//			m.lastKey = "←"
//		case key.Matches(msg, m.keys.Right):
//			m.lastKey = "→"
//		case key.Matches(msg, m.keys.Help):
//			m.help.ShowAll = !m.help.ShowAll
//		case key.Matches(msg, m.keys.Quit):
//			m.quitting = true
//			return m, tea.Quit
//		}
//	}
//
//	return m, nil
//}
//
//func (m helpRegion) View() string {
//	if m.quitting {
//		return "Bye!\n"
//	}
//
//	var status string
//	if m.lastKey == "" {
//		status = "Waiting for input..."
//	} else {
//		status = "You chose: " + m.inputStyle.Render(m.lastKey)
//	}
//
//	helpView := m.help.View(m.keys)
//	height := 8 - strings.Count(status, "\n") - strings.Count(helpView, "\n")
//
//	return "\n" + status + strings.Repeat("\n", height) + helpView
//}

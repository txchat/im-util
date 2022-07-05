package tui

//
//import (
//	"fmt"
//	"github.com/charmbracelet/bubbles/textarea"
//	"github.com/charmbracelet/bubbles/viewport"
//	tea "github.com/charmbracelet/bubbletea"
//	"github.com/charmbracelet/lipgloss"
//	"strings"
//)
//
//type sessionPage struct {
//	//显示聊天记录区域
//	viewport    viewport.Model
//
//	//发送区域
//
//
//
//	messages    []string
//	textarea    textarea.Model
//	senderStyle lipgloss.Style
//	err         error
//}
//
//func setViewPort() {
//	vp := viewport.New(30, 10)
//	vp.SetContent(``)
//}
//
//func NewSessionPage(chType int32, target string) sessionPage {
//	ta := textarea.New()
//	ta.Placeholder = "Send a message..."
//	ta.Focus()
//
//	ta.Prompt = "┃ "
//	ta.CharLimit = 280
//
//	ta.SetWidth(30)
//	ta.SetHeight(3)
//
//	// Remove cursor line styling
//	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
//
//	ta.ShowLineNumbers = false
//
//
//	ta.KeyMap.InsertNewline.SetEnabled(false)
//
//	return sessionPage{
//		textarea:    ta,
//		messages:    []string{},
//		viewport:    vp,
//		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
//		err:         nil,
//	}
//}
//
//func (m sessionPage) Init() tea.Cmd {
//	return textarea.Blink
//}
//
//func (m sessionPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//	var (
//		tiCmd tea.Cmd
//		vpCmd tea.Cmd
//	)
//
//	m.textarea, tiCmd = m.textarea.Update(msg)
//	m.viewport, vpCmd = m.viewport.Update(msg)
//
//	switch msg := msg.(type) {
//	case tea.KeyMsg:
//		switch msg.Type {
//		case tea.KeyEsc:
//			fmt.Println(m.textarea.Value())
//			return m, tea.Quit
//		case tea.KeyEnter:
//			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
//			m.viewport.SetContent(strings.Join(m.messages, "\n"))
//			m.textarea.Reset()
//			m.viewport.GotoBottom()
//		}
//
//	// We handle errors just like any other message
//	case errMsg:
//		m.err = msg
//		return m, nil
//	}
//
//	return m, tea.Batch(tiCmd, vpCmd)
//}
//
//func (m sessionPage) View() string {
//	return fmt.Sprintf(
//		"%s\n\n%s",
//		m.viewport.View(),
//		m.textarea.View(),
//	) + "\n\n"
//}

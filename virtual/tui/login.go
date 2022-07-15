package tui

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/txchat/im-util/internal/device"
	"github.com/txchat/im-util/internal/user"
	"github.com/txchat/im-util/protocol/wallet"
	"github.com/txchat/imparse/proto"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type LoginPage struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode textinput.CursorMode
}

func NewLoginPage() LoginPage {
	m := LoginPage{
		inputs: make([]textinput.Model, 6),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 64

		switch i {
		case 0:
			t.Placeholder = "应用ID"
			t.SetValue("dtalk")
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "设备类型"
			t.SetValue("Android")
		case 2:
			t.Placeholder = "设备名称"
			t.SetValue("虚拟驱动")
		case 3:
			t.Placeholder = "设备唯一识别号"
			t.SetValue("3ade6a21-a0d7-48ce-94a2-2f3567adc468")
		case 4:
			t.Placeholder = "服务端地址"
			t.SetValue("localhost:3302")
		case 5:
			t.Placeholder = "助记词"
			t.SetValue("")
		}

		m.inputs[i] = t
	}

	return m
}

func (m LoginPage) Init() tea.Cmd {
	return textinput.Blink
}

func (m LoginPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > textinput.CursorHide {
				m.cursorMode = textinput.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].SetCursorMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				// login
				mp := NewMainPage()
				return mp, mp.Init()
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *LoginPage) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *LoginPage) doLogin(mnc, uuid, deviceName string, deviceType int32, appId, server string) error {
	// check
	w, err := wallet.NewWalletFromMnemonic(mnc)
	if err != nil {
		return err
	}
	md, err := wallet.FormatMetadataFromWallet(0, w)
	if err != nil {
		return err
	}
	u := user.NewUser(md.GetAddress(), md.GetPrivateKey(), md.GetPublicKey())
	d := device.NewDevice(uuid, deviceName, proto.Device(deviceType), zerolog.Logger{}, u)
	err = d.DialIMServer(appId, server, nil)
	if err != nil {
		return err
	}
	err = d.TurnOn()
	if err != nil {
		return err
	}
	return nil
}

func (m LoginPage) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}

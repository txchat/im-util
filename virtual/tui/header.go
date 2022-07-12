package tui

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"text/template"

	tea "github.com/charmbracelet/bubbletea"
)

type headerRegion struct {
	content []byte
}

func NewHeaderRegion() headerRegion {
	h := headerRegion{
		content: make([]byte, 0),
	}
	if err := h.loadHeadContent(); err != nil {
		panic(err)
	}
	return h
}

func (m headerRegion) Init() tea.Cmd {
	return nil
}

func (m headerRegion) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m headerRegion) View() string {
	return string(m.content)
}

func (m *headerRegion) loadHeadContent() error {
	//init head
	h := HeadInfo{
		DeviceName: "",
		DeviceType: "",
		UUID:       "",
		APPID:      "",
		UID:        "",
		Server:     "",
	}

	headTmpl, err := template.New("head").Parse(HeadText)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	wr := bufio.NewWriter(buf)
	err = headTmpl.Execute(wr, h)
	if err != nil {
		return err
	}
	err = wr.Flush()
	if err != nil {
		return err
	}
	m.content, err = ioutil.ReadAll(buf)
	if err != nil {
		return err
	}
	return nil
}

type HeadInfo struct {
	DeviceName string
	DeviceType string
	UUID       string
	APPID      string
	UID        string
	Server     string
}

var HeadText = `
设备名称:{{.DeviceName}}
设备类型:{{.DeviceType}}
设备唯一识别号:{{.UUID}}

应用名称:{{.APPID}}
用户ID:{{.UID}}
连接服务器地址:{{.Server}}
----------------------------------------------------------
`

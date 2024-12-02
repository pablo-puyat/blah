package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error
type logMsg string  // Changed from testMsg to logMsg as this is a real message type

type model struct {
	lines    []string
	spinner  spinner.Model
	quitting bool
	err      error
	logChan  chan string
	style    lipgloss.Style
}

var quitKeys = key.NewBinding(
	key.WithKeys("q", "esc", "ctrl+c"),
	key.WithHelp("", "press q to quit"),
)

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		lines:   make([]string, 0),
		spinner: s,
		logChan: make(chan string),
		style:   lipgloss.NewStyle().Padding(0, 1),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.waitForLog,
	)
}

func (m model) waitForLog() tea.Msg {
	msg := <-m.logChan
	return logMsg(msg)  // Using logMsg instead of testMsg
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) {
			m.quitting = true
			return m, tea.Quit
		}
		return m, nil

	case logMsg:  // Changed from testMsg to logMsg
		m.lines = append(m.lines, string(msg))
		return m, m.waitForLog

	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("\n %s Logging...\n\n", m.spinner.View()))

	for _, line := range m.lines {
		b.WriteString(m.style.Render(line))
		b.WriteString("\n")
	}

	b.WriteString(fmt.Sprintf("\n %s\n", quitKeys.Help().Desc))

	if m.quitting {
		b.WriteString("\n")
	}

	return b.String()
}

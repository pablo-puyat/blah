package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type ErrMsg error

type Model struct {
	lines    []string
	quitting bool
	err      error
	sub      <-chan string
}

var quitKeys = key.NewBinding(
	key.WithKeys("q", "esc", "ctrl+c"),
	key.WithHelp("", "press q to quit"),
)

func waitForLog(sub <-chan string) tea.Cmd {
	return func() tea.Msg {
		return string(<-sub)
	}
}

func NewModel(sub <-chan string) Model {
	return Model {
		lines: []string{"Hello","Bye"},
		quitting: false,
		err: nil,
		sub: sub,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		waitForLog(m.sub),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) {
			m.quitting = true
			return m, tea.Quit
		}
		return m, nil

	case string:
		m.lines = append(m.lines, string(msg))
		return m, waitForLog(m.sub)

	case ErrMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		return m, cmd
	}
}

func (m Model) View() string {
	var s string
    for _, item := range m.lines {
        s += fmt.Sprintf("%s\n", item)
    }
    return s
}

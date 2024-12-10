package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	lines  []string
	source string
	ready  bool
}

func New() *model {
	return &model{ready: false, lines: []string{"Watching"}}
}
func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case string:
		m.lines = append(m.lines, msg)
	case tea.WindowSizeMsg:
		m.ready = true
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "Initializing"
	} 
	return strings.Join(m.lines, "\n")
}

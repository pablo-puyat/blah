package tui

import (
	"fmt"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestModelHandlesLogMessage(t *testing.T) {
	m := initialModel()
	p := tea.NewProgram(m, tea.WithoutRenderer())
	
	go func() {
		if _, err := p.Run(); err != nil {
			t.Errorf("Failed to run program: %v", err)
		}
	}()

	// Wait for program to initialize
	time.Sleep(100 * time.Millisecond)

	// Send a log message through the channel, as it would happen in real usage
	m.logChan <- "Test log entry 123"

	time.Sleep(100 * time.Millisecond)

	// Test quit functionality
	p.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	time.Sleep(100 * time.Millisecond)
	
	p.Kill()
}

func TestModelHandlesError(t *testing.T) {
	m := initialModel()
	p := tea.NewProgram(m, tea.WithoutRenderer())

	go func() {
		if _, err := p.Run(); err != nil {
			t.Errorf("Failed to run program: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	p.Send(errMsg(fmt.Errorf("test error")))
	time.Sleep(100 * time.Millisecond)
	
	p.Kill()
}

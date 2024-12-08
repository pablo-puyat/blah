package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"strings"
	"testing"
	"time"
	"github.com/pablo-puyat/blahblah/internal/testutil"
)

// Test function
func TestComponent(t *testing.T) {
	// Initialize your model
	m := NewModel()

	// Simulate initial window size message and type assert the result
	updatedModel, _ := m.Update(tea.WindowSizeMsg{
		Width:  80,
		Height: 24,
	})

	// Type assert the interface back to our concrete type
	m, ok := updatedModel.(Model)
	if !ok {
		t.Fatal("Could not type assert model")
	}

	if !m.ready {
		t.Error("Component should be ready after receiving WindowSizeMsg")
	}

	if !strings.Contains(m.View(), "watching") {
		t.Error("Viewport should contain watching")
	}
}

func TestLogMessagesAppearInOutput(t *testing.T) {
	m := NewModel()
	p := tea.NewProgram(m)

	go func() {
		if _, err := p.Run(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}()

	updatedModel, _ := m.Update(tea.WindowSizeMsg{
		Width:  80,
		Height: 24,
	})
	m, ok := updatedModel.(Model)
	if !ok {
		t.Fatal("Could not type assert model")
	}

	time.Sleep(300 * time.Millisecond)

	if !m.ready {
		t.Error("Component should be ready after receiving WindowSizeMsg")
	}

/*
	updatedModel, _ = m.Update("sdfsdfsdf")
	m, ok = updatedModel.(Model)
	if !ok {
		t.Fatal("Could not type assert model")
	}
*/
	l := "test log message"

	p.Send(l)
	time.Sleep(300 * time.Millisecond)
	if !strings.Contains(m.View(), l) {
		t.Error("output should contain " + l + "\n")
		testutil.Logger.Printf("TUI test")
	}
	testutil.Logger.Printf("end")
	p.Quit()
}

package tui

import (
	"bytes"
	"strings"
	tea "github.com/charmbracelet/bubbletea"
	"testing"
	"time"
)

func TestLogMessagesAppearInOutput(t *testing.T) {
	var buf bytes.Buffer
	// Initialize program with test renderer
	ch := make(chan string)
	defer close(ch)
	p := tea.NewProgram(NewModel(ch), tea.WithOutput(&buf))

	// Run program in background
	go func() {
		if _, err := p.Run(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}()

	// Wait for initialization
	time.Sleep(100 * time.Millisecond)

	// Send log message via channel
	expectedLog := "test log message"
	ch <- expectedLog

	// Wait for processing
	time.Sleep(100 * time.Millisecond)

    if !strings.Contains(buf.String(), expectedLog) {
        t.Error("output should contain 'expected'")
    }  
	time.Sleep(100 * time.Millisecond)
	// Clean up
	p.Quit()
}

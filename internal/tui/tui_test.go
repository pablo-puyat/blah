package tui

import (
	"crypto/rand"
	"encoding/base64"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
	"testing"
	"time"
)

func TestLogMessagesAppearInOutput(t *testing.T) {
	m := New()
	p := tea.NewProgram(m)
	go func() {
		if _, err := p.Run(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}()
	if !strings.Contains(m.View(), "Initializing") {
		t.Error("output should contain Initializing", m.View())
	}

	time.Sleep(100 * time.Millisecond)

	updatedModel, _ := m.Update(tea.WindowSizeMsg{
		Width:  80,
		Height: 24,
	})

	// Type assert the interface back to our concrete type
	m, ok := updatedModel.(*model)
	if !ok {
		t.Fatal("Could not type assert model")
	}

	if !m.ready {
		t.Error("Component should be ready after receiving WindowSizeMsg")
	}
	for i := 0; i < 4; i++ {
		rm , err := generateRandomString(20)
		if err != nil {
			t.Fatal("unable to generate random strings")
		}

		p.Send(rm)
		time.Sleep(300 * time.Millisecond)
		if !strings.Contains(m.View(), rm) {
			t.Error("output should contain "+ rm +"\n", m.View())
		}

	}
	p.Quit()
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

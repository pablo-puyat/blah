package filewatcher

import (
	"os"
	"testing"
	"time"
)

func TestWatcherOnlyReadsNewLines(t *testing.T) {
	// Create temp file with initial content that should be ignored
	tmpfile, err := os.CreateTemp("", "test*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	// Write initial content that should be ignored by watcher
	initialContent := "old line 1\nold line 2\n"
	if err := os.WriteFile(tmpfile.Name(), []byte(initialContent), 0644); err != nil {
		t.Fatalf("Failed to write initial content: %v", err)
	}

	// Start the watcher
	fw, err := NewFileWatcher(tmpfile.Name())
	if err != nil {
		
	}
	done := make(chan bool)
	go fw.Watch(done)
	defer close(done)

	// Give the watcher a moment to start up
	time.Sleep(100 * time.Millisecond)

	// Append new entries to the file
	f, err := os.OpenFile(tmpfile.Name(), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("Failed to open file for appending: %v", err)
	}

	newEntries := []string{"new line 1", "new line 2"}
	for _, entry := range newEntries {
		if _, err := f.WriteString(entry + "\n"); err != nil {
			t.Fatalf("Failed to append to file: %v", err)
		}
	}
	f.Close()

	// Verify only new lines are detected
	for i, exp := range newEntries {
		select {
		case line := <-fw.Lines:
			if line != exp {
				t.Errorf("Line %d: expected %q, got %q", i+1, exp, line)
			}
		case <-time.After(time.Second):
			t.Fatalf("Timeout waiting for line %d", i+1)
		}
	}

	// Verify no old content was sent
	select {
	case line := <-fw.Lines:
		t.Errorf("Unexpected line received: %q (watcher should ignore existing content)", line)
	case <-time.After(100 * time.Millisecond):
		// This is good - we expect no more lines
	}
}

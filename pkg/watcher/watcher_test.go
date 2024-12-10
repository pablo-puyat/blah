package watcher

import (
	"fmt"
	"os"
	"testing"
	"time"
	
)

func TestWatcherOnlyReadsNewLines(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	initialContent := "old line 1\nold line 2\n"
	if _, err := tmpfile.WriteString(initialContent); err != nil {
		t.Fatalf("Failed to write initial content: %v", err)
	}
	info, err := os.Stat(tmpfile.Name())
	if err != nil {
		t.Logf("Error getting file info: %v", err)
	} else {
		t.Logf("File size: %d", info.Size())
	}

	fw, err := New(tmpfile.Name())
	if err != nil {
		t.Fatalf("Error instantiating file watcher")
	}
	lines := make(chan *Line, 10)

	if err := fw.Watch(lines); err != nil {
		t.Fatalf("Unable to open file to watch")
	}

	// Give the watcher a moment to start up
	time.Sleep(300 * time.Millisecond)

	for i := 0; i < 5; i++ {
		if _, err := tmpfile.WriteString(fmt.Sprintf("test line %d\n", i)); err != nil {
			t.Fatalf("Failed to append to file: %v", err)
		}
	}

	done := make(chan bool)

	go func() {
		for i := 0; i < 5; i++ {
			line := <-lines
			if line.Content != fmt.Sprintf("test line %d", i) {
				t.Errorf("Expected: %s\nReceived: ", line)
			}
		}
		done <- true
	}()
	<-done
}

package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pablo-puyat/blahblah/internal/filewatcher"
	"github.com/pablo-puyat/blahblah/internal/tui"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: blah <logfile>")
		os.Exit(1)
	}

	// Initialize file watcher
	fw, err := filewatcher.NewFileWatcher(os.Args[1])
	if err != nil {
		fmt.Printf("Error initializing file watcher: %v\n", err)
		os.Exit(1)
	}

	// Start watching file in background
	done := make(chan bool)

	// Create and start TUI
	p := tea.NewProgram(tui.NewModel())

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}

	go func() {
		if err := fw.Watch(done); err != nil {
			fmt.Printf("Watcher error: %v\n", err)
			os.Exit(1)
		}
		for {
			msg, more := <-fw.Lines
			if more {
				p.Send(msg)
			} else {
				done <- true
				return
			}
		}
	}()
	<-done
	close(done)
}

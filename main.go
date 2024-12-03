package main

import (
	"fmt"
	"os"

	"github.com/pablo-puyat/blahblah/internal/filewatcher"
	"github.com/pablo-puyat/blahblah/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: blahblah <logfile>")
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
	go func() {
		if err := fw.Watch(done); err != nil {
			fmt.Printf("Watcher error: %v\n", err)
			os.Exit(1)
		}
	}()

	// Create and start TUI
	program := tea.NewProgram(tui.NewModel(fw.Lines))

	// Run the program
	if _, err := program.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}

	// Cleanup
	close(done)
}

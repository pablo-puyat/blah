package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"blah/internal/watcher"
	"blah/internal/tui"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: blah <logfile>")
		os.Exit(1)
	}

	fw, err := watcher.New(os.Args[1])
	if err != nil {
		fmt.Printf("Error initializing file watcher: %v\n", err)
		os.Exit(1)
	}

	done := make(chan bool)
	lines := make(chan *watcher.Line, 50)

	p := tea.NewProgram(tui.New(), tea.WithAltScreen())
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Printf("Unable to open log file: %v", err)
		os.Exit(1)
	}
	defer f.Close()

	if err := fw.Watch(lines); err != nil {
		log.Printf("Unable to start watcher: %v", err)
		os.Exit(1)
	}
	go func() {
		if _, err := p.Run(); err != nil {
			log.Printf("Error running program: %v", err)
			os.Exit(1)
		}
	}()

	go func() {
		for {
			l, more := <-lines
			if more {
				p.Send(l)
			} else {
				done <- true
				return
			}
		}
	}()

	<-done
	close(done)
	p.Quit()
}

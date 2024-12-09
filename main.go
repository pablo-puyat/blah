package main

import (
	"fmt"
	"log"
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

	fw, err := filewatcher.NewFileWatcher(os.Args[1])
	if err != nil {
		fmt.Printf("Error initializing file watcher: %v\n", err)
		os.Exit(1)
	}

	done := make(chan bool)
	lines := make(chan string)

	p := tea.NewProgram(tui.New(), tea.WithAltScreen())
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal("Unable to open log file", err)
	}
	defer f.Close()
	log.Print("About to launch a thousand ships")
	go func() {
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running program: %v\n", err)
			os.Exit(1)
		}
	}()

	go func() {
		log.Print("Watch start")
		if err := fw.Watch(lines); err != nil {
			fmt.Printf("Watcher error: %v\n", err)
			os.Exit(1)
		}
		log.Print("Watch end")
	}()

	go func() {
		log.Print("Selector started")
		for {
			l, more := <-lines
			if more {
				log.Print("Received job", l)
				p.Send(l)
			} else {
				log.Print("No more logs")
				done<- true
				return
			}
		}
	}()

	<-done
	close(done)
	log.Print("Just waiting")
	p.Quit()
}

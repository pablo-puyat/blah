package watcher

import (
	"bufio"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

type Watcher struct {
	filename string
	writer   chan string
}

type Line struct {
	Content   string
	Source    string
	Timestamp time.Time
}

func New(filename string) (*Watcher, error) {
	return &Watcher{
		filename: filename,
	}, nil
}

func (fw *Watcher) Watch(writer chan<- *Line) error {
	file, err := os.Open(fw.filename)
	if err != nil {
		return err
	}

	go func() {
		defer file.Close()
		var mark int64
		first := true
		for {
			info, err := file.Stat()
			if err != nil {
				log.Fatalf("error getting file size: %v", err)
			}
			if first {
				mark = info.Size()
				first = false
			}
			if mark < info.Size() {
				file.Seek(mark, 0)
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					writer <- &Line{
						Content:   scanner.Text(),
						Source:    fw.filename,
						Timestamp: time.Now(),
					}
				}
			}
			mark = info.Size()
			time.Sleep(6 * time.Millisecond)
		}
	}()

	return nil
}

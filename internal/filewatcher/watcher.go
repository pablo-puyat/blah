package filewatcher

import (
	"bufio"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

type FileWatcher struct {
	filename string
	writer   chan string
}

type Line struct {
	Content   string
	Source    string
	Timestamp time.Time
}

func NewFileWatcher(filename string) (*FileWatcher, error) {
	return &FileWatcher{
		filename: filename,
	}, nil
}

func (fw *FileWatcher) Watch(writer chan<- *Line) error {
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

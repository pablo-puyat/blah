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

func NewFileWatcher(filename string) (*FileWatcher, error) {
	return &FileWatcher{
		filename: filename,
	}, nil
}

func (fw *FileWatcher) Watch(writer chan<- string) error {
	file, err := os.Open(fw.filename)
	if err != nil {
		return err
	}

	go func() {
		log.Print("In go routine")
		defer file.Close()
		var fileMark int64
		firstRun := true
		for {
				info, err := file.Stat()
				if err != nil {
					log.Fatalf("error getting file size: %v", err)
				}
				if firstRun {
					fileMark = info.Size()
					firstRun = false
				}
				if fileMark < info.Size() {
					file.Seek(fileMark, 0)
					scanner := bufio.NewScanner(file)
					for scanner.Scan() {
						writer <- scanner.Text()
						log.Print("Wrote message")
					}
				}
				fileMark = info.Size()
				time.Sleep(time.Second)
			}
	}()

	return nil
}

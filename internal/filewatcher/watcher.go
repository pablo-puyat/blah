package filewatcher

import (
    "bufio"
    "fmt"
    "os"
    "time"
)

type FileWatcher struct {
    filename     string
    Lines        chan string
    file         *os.File
    size         int64
}

func NewFileWatcher(filename string) (*FileWatcher, error) {
    file, err := os.Open(filename)
    if err !=nil {
        return nil, fmt.Errorf("error opening file: %v", err)
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, fmt.Errorf("error getting file size")
    }
    print(info.Size()) 
    return &FileWatcher{
        filename:     filename,
        file:         file,
        Lines:        make(chan string),
        size: info.Size(),
    }, nil
}

func (fw *FileWatcher) Watch(done chan bool) error {
    //defer fw.file.Close()
    //defer close(fw.Lines)

    for {
        info, err := fw.file.Stat()
        if err != nil {
            return fmt.Errorf("error getting file size: %v", err)
        }
        print(info.Size())
        if info.Size() > fw.size {
            for {
                select {
                case <-done:
                    return nil
                default:
                    scanner := bufio.NewScanner(fw.file)
                    fw.file.Seek(fw.size, 0)
                    for scanner.Scan() {
                        fw.Lines <- scanner.Text()
                    }
                    fw.file.Close()
                    time.Sleep(100 * time.Millisecond)
                }
            }
        }
    }
}

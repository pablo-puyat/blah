package testutil

import (
    "log"
    "os"
    "path/filepath"
)

var Logger *log.Logger

func init() {
    // Create logs directory if it doesn't exist
    err := os.MkdirAll("../../logs", 0755)
    if err != nil {
        log.Fatal(err)
    }

    // Create or open log file in the logs directory
    f, err := os.OpenFile(filepath.Join("../../logs", "test.log"), 
        os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal(err)
    }

    Logger = log.New(f, "", log.Ltime|log.Lshortfile)
}

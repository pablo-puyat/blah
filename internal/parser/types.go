package parser

import "time"

type LogLevel string

const (
    LevelDebug   LogLevel = "DEBUG"
    LevelInfo    LogLevel = "INFO"
    LevelWarning LogLevel = "WARNING"
    LevelError   LogLevel = "ERROR"
    LevelFatal   LogLevel = "FATAL"
)

type StackFrame struct {
    File       string
    Line       int
    Function   string
    Args       string
    Raw        string
}

type LogEntry struct {
    Timestamp  time.Time
    Level      LogLevel
    Message    string
    Channel    string
    Context    map[string]any
    Exception  string
    StackTrace []StackFrame
    Raw        string
}

type Parser interface {
    Parse(lines []string) (*LogEntry, error)
    RequiresMultiline() bool
    IsStartOfEntry(line string) bool
}

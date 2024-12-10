package laravel

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"blah/internal/parser"
)

type Parser struct {
    timeFormat string
    timezone   *time.Location
}

func New(opts ...parser.Option) (parser.Parser, error) {
    cfg := &parser.Config{
        TimeFormat: "2006-01-02 15:04:05",
        TimeZone:   "UTC",
    }
    
    for _, opt := range opts {
        opt(cfg)
    }
    
    tz, err := time.LoadLocation(cfg.TimeZone)
    if err != nil {
        return nil, fmt.Errorf("invalid timezone: %w", err)
    }
    
    return &Parser{
        timeFormat: cfg.TimeFormat,
        timezone:   tz,
    }, nil
}

func (p *Parser) RequiresMultiline() bool {
    return true
}

func (p *Parser) IsStartOfEntry(line string) bool {
    return logLineRegex.MatchString(line)
}

func (p *Parser) Parse(lines []string) (*parser.LogEntry, error) {
    if len(lines) == 0 {
        return nil, fmt.Errorf("no lines to parse")
    }
    
    matches := logLineRegex.FindStringSubmatch(lines[0])
    if matches == nil {
        return nil, fmt.Errorf("invalid log format")
    }
    
    timestamp, err := time.ParseInLocation(p.timeFormat, matches[1], p.timezone)
    if err != nil {
        return nil, fmt.Errorf("invalid timestamp: %w", err)
    }
    
    // Initialize the log entry with basic info
    entry := &parser.LogEntry{
        Timestamp: timestamp,
        Channel:   matches[2],
        Level:     parser.LogLevel(strings.ToUpper(matches[3])),
        Message:   matches[4],
        Raw:      strings.Join(lines, "\n"),
    }

    // If we have multiple lines, handle the context and stack trace
    if len(lines) > 1 {
        // Find JSON context in the first line if it exists
        if jsonMatch := jsonContextRegex.FindString(lines[0]); jsonMatch != "" {
            var contextData struct {
                Exception  string         `json:"exception"`
                Context   map[string]any `json:"context"`
            }
            if err := json.Unmarshal([]byte(jsonMatch), &contextData); err == nil {
                entry.Exception = contextData.Exception
                entry.Context = contextData.Context
            }
        }

        // Look for stack trace in remaining lines
        var frames []parser.StackFrame
        inStackTrace := false
        
        for _, line := range lines[1:] {
            // Check if we're entering the stack trace section
            if strings.Contains(line, "[stacktrace]") {
                inStackTrace = true
                continue
            }
            
            // If we're in the stack trace section, parse frames
            if inStackTrace {
                if matches := stackFrameRegex.FindStringSubmatch(line); matches != nil {
                    lineNum, _ := strconv.Atoi(matches[3])
                    frames = append(frames, parser.StackFrame{
                        File:     matches[2],
                        Line:     lineNum,
                        Function: matches[4],
                        Raw:      line,
                    })
                }
            }
        }
        
        if len(frames) > 0 {
            entry.StackTrace = frames
        }
    }
    
    return entry, nil
}

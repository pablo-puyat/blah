package laravel

import (
    "reflect"
    "testing"
    "time"

    "github.com/davecgh/go-spew/spew"
    "blah/pkg/parser"
)

func TestParser_Parse(t *testing.T) {
    // Create parser with fixed timezone for consistent tests
    p, err := New()
    if err != nil {
        t.Fatalf("Failed to create parser: %v", err)
    }

    tests := []struct {
        name    string
        lines   []string
        want    *parser.LogEntry
        wantErr bool
    }{
        {
            name: "simple info log",
            lines: []string{
                "[2024-11-30 00:53:41] local.INFO: Hello",
            },
            want: &parser.LogEntry{
                Timestamp: mustParseTime("2024-11-30 00:53:41"),
                Channel:   "local",
                Level:    parser.LevelInfo,
                Message:  "Hello",
                Raw:     "[2024-11-30 00:53:41] local.INFO: Hello",
            },
            wantErr: false,
        },
        {
            name: "error with context",
            lines: []string{
                `[2024-11-30 00:53:41] local.ERROR: Call to a member function createToken() on null {"exception":"Error message"}`,
            },
            want: &parser.LogEntry{
                Timestamp: mustParseTime("2024-11-30 00:53:41"),
                Channel:   "local",
                Level:    parser.LevelError,
                Message:  `Call to a member function createToken() on null`,
                Raw:      `[2024-11-30 00:53:41] local.ERROR: Call to a member function createToken() on null {"exception":"Error message"}`,
            },
            wantErr: false,
        },
        {
            name: "error with stack trace",
            lines: []string{
                `[2024-11-30 00:53:41] local.ERROR: Call to a member function createToken() on null {"exception":"[object] (Error(code: 0): Call to a member function createToken() on null at /routes/api.php:15)"}`,
                "[stacktrace]",
                "#0 /vendor/laravel/framework/src/Illuminate/Routing/CallableDispatcher.php(40): RouteFileRegistrar->{closure}(Object)",
                "#1 /vendor/laravel/framework/src/Illuminate/Routing/Route.php(240): CallableDispatcher->dispatch(Object)",
            },
            want: &parser.LogEntry{
                Timestamp: mustParseTime("2024-11-30 00:53:41"),
                Channel:   "local",
                Level:    parser.LevelError,
                Message:  "Call to a member function createToken() on null",
                Exception: "[object] (Error(code: 0): Call to a member function createToken() on null at /routes/api.php:15)",
                StackTrace: []parser.StackFrame{
                    {
                        File:     "/vendor/laravel/framework/src/Illuminate/Routing/CallableDispatcher.php",
                        Line:     40,
                        Function: "RouteFileRegistrar->{closure}(Object)",
                        Raw:      "#0 /vendor/laravel/framework/src/Illuminate/Routing/CallableDispatcher.php(40): RouteFileRegistrar->{closure}(Object)",
                    },
                    {
                        File:     "/vendor/laravel/framework/src/Illuminate/Routing/Route.php",
                        Line:     240,
                        Function: "CallableDispatcher->dispatch(Object)",
                        Raw:      "#1 /vendor/laravel/framework/src/Illuminate/Routing/Route.php(240): CallableDispatcher->dispatch(Object)",
                    },
                },
                Raw: `[2024-11-30 00:53:41] local.ERROR: Call to a member function createToken() on null {"exception":"[object] (Error(code: 0): Call to a member function createToken() on null at /routes/api.php:15)"}
[stacktrace]
#0 /vendor/laravel/framework/src/Illuminate/Routing/CallableDispatcher.php(40): RouteFileRegistrar->{closure}(Object)
#1 /vendor/laravel/framework/src/Illuminate/Routing/Route.php(240): CallableDispatcher->dispatch(Object)`,
            },
            wantErr: false,
        },
        {
            name: "invalid log format",
            lines: []string{
                "This is not a Laravel log",
            },
            want:    nil,
            wantErr: true,
        },
        {
            name:    "empty input",
            lines:   []string{},
            want:    nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := p.Parse(tt.lines)
            if (err != nil) != tt.wantErr {
                t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if tt.wantErr {
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("Got:\n%s\n\nWant\n%s\n\n", spew.Sdump(got), spew.Sdump(tt.want))
            }
        })
    }
}

func mustParseTime(value string) time.Time {
    t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.UTC)
    if err != nil {
        panic(err)
    }
    return t
}
/*
func TestParser_RequiresMultiline(t *testing.T) {
    p, err := New()
    if err != nil {
        t.Fatalf("Failed to create parser: %v", err)
    }
    
    if !p.RequiresMultiline() {
        t.Error("RequiresMultiline() = false, want true")
    }
}

func TestParser_IsStartOfEntry(t *testing.T) {
    p, err := New()
    if err != nil {
        t.Fatalf("Failed to create parser: %v", err)
    }

    tests := []struct {
        name  string
        line  string
        want  bool
    }{
        {
            name: "valid log start",
            line: "[2024-11-30 00:53:41] local.INFO: Message",
            want: true,
        },
        {
            name: "stack trace line",
            line: "#0 /path/to/file.php(123): Method()",
            want: false,
        },
        {
            name: "invalid format",
            line: "Not a log line",
            want: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := p.IsStartOfEntry(tt.line); got != tt.want {
                t.Errorf("IsStartOfEntry() = %v, want %v", got, tt.want)
            }
        })
    }
}
*/

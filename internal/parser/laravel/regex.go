package laravel

import "regexp"

var (
    // [2024-11-30 00:53:41] local.INFO: Message
    logLineRegex = regexp.MustCompile(`^\[(.+?)\] (\w+)\.(\w+): (.+?)( \{.*\})?$`)
    
    // Matches the JSON part that contains exception and context
    jsonContextRegex = regexp.MustCompile(`\{.*\}$`)
    
    // #0 /path/to/file.php(123): Class->method(args)
    stackFrameRegex = regexp.MustCompile(`^#(\d+)\s+(.+?)\((\d+)\): (.+)$`)
)

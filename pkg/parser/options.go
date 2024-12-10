package parser

type Config struct {
    TimeFormat string
    TimeZone   string
}

type Option func(*Config)

func WithTimeFormat(format string) Option {
    return func(c *Config) {
        c.TimeFormat = format
    }
}

func WithTimeZone(tz string) Option {
    return func(c *Config) {
        c.TimeZone = tz
    }
}

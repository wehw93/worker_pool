// standard.go
package logger

import (
    "log"
    "os"
)


type StandardLogger struct {
    logger *log.Logger
    level  LogLevel
}

func NewStandardLogger(config Config) *StandardLogger {
    return &StandardLogger{
        logger: log.New(os.Stdout, config.Prefix, log.LstdFlags),
        level:  config.Level,
    }
}

func (l *StandardLogger) Debug(msg string) {
    if l.level <= LevelDebug {
        l.logger.Printf("[DEBUG] %s", msg)
    }
}

func (l *StandardLogger) Info(msg string) {
    if l.level <= LevelInfo {
        l.logger.Printf("[INFO] %s", msg)
    }
}

func (l *StandardLogger) Error(msg string) {
    if l.level <= LevelError {
        l.logger.Printf("[ERROR] %s", msg)
    }
}

func (l *StandardLogger) Printf(format string, v ...interface{}) {
    l.logger.Printf(format, v...)
}
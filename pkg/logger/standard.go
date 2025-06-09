
package logger

import (
	"log"
	"os"
)

type StandardLogger struct {
	config Config
	logger *log.Logger
}

func NewStandardLogger(config Config) Logger {
	return &StandardLogger{
		config: config,
		logger: log.New(os.Stdout, config.Prefix, log.LstdFlags),
	}
}

func (l *StandardLogger) Debug(msg string) {
	if l.config.Level <= LevelDebug {
		l.logger.Printf("[DEBUG] %s", msg)
	}
}

func (l *StandardLogger) Info(msg string) {
	if l.config.Level <= LevelInfo {
		l.logger.Printf("[INFO] %s", msg)
	}
}

func (l *StandardLogger) Error(msg string) {
	if l.config.Level <= LevelError {
		l.logger.Printf("[ERROR] %s", msg)
	}
}

func (l *StandardLogger) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

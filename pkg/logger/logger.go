// logger.go
package logger


type Logger interface {
    Debug(msg string)
    Info(msg string)
    Error(msg string)
    Printf(format string, v ...interface{})
}


type LogLevel int

const (
    LevelDebug LogLevel = iota
    LevelInfo
    LevelError
)


type Config struct {
    Level  LogLevel
    Prefix string
}
package logger

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
)

const (
	LogLevelDebug   = iota // 0
	LogLevelInfo           // 1
	LogLevelWarning        // 2
	LogLevelError          // 3
	LogLevelFatal
	LogLevelPanic
	LogLevelNone
	LogLevelDisabled
)

var LogLevelMap = map[int]string{
	LogLevelDebug:    "DEBUG",
	LogLevelInfo:     "INFO",
	LogLevelWarning:  "WARNING",
	LogLevelError:    "ERROR",
	LogLevelFatal:    "FATAL",
	LogLevelPanic:    "PANIC",
	LogLevelNone:     "NONE",
	LogLevelDisabled: "Disabled",
}

type Logger struct {
	logger       *log.Logger
	currentLevel int
}

var (
	// DefaultLoggerFlags default flags for logger
	DefaultLoggerFlags = log.LstdFlags | log.Lmsgprefix
)

type Interface interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

var _ Interface = (*Logger)(nil)

func New(prefix string, logLevel string) *Logger {
	l := Logger{
		logger: log.New(os.Stdout, prefix, DefaultLoggerFlags),
	}
	l.SetLogLevel(logLevel)
	return &l
}

func (l *Logger) SetLogLevel(level string) {
	switch strings.ToLower(level) {
	case "error":
		l.currentLevel = LogLevelError
	case "warn":
		l.currentLevel = LogLevelWarning
	case "info":
		l.currentLevel = LogLevelInfo
	case "debug":
		l.currentLevel = LogLevelDebug
	default:
		l.currentLevel = LogLevelInfo
	}
}

func (l *Logger) logWithCaller(level int, format string, v ...interface{}) {
	if l.currentLevel <= level {
		_, file, line, _ := runtime.Caller(2)
		filename := path.Base(file)
		callerInfo := fmt.Sprintf("[%s:%d]", filename, line)
		logLevelInfo := fmt.Sprintf("[%s]", LogLevelMap[level])
		l.logger.Println(logLevelInfo, callerInfo, fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Debug(v ...interface{}) {
	l.logWithCaller(LogLevelDebug, fmt.Sprint(v...))
}

func (l *Logger) Info(v ...interface{}) {
	l.logWithCaller(LogLevelInfo, fmt.Sprint(v...))
}

func (l *Logger) Warn(v ...interface{}) {
	l.logWithCaller(LogLevelWarning, fmt.Sprint(v...))
}

func (l *Logger) Error(v ...interface{}) {
	l.logWithCaller(LogLevelError, fmt.Sprint(v...))
}

func (l *Logger) Fatal(v ...interface{}) {
	l.logWithCaller(LogLevelError, fmt.Sprint(v...))
	os.Exit(1)
}

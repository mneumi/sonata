package log

import (
	"fmt"
	"io"
	"os"
)

type LoggerLevel int

const (
	LevelDebug LoggerLevel = iota
	LevelInfo
	LevelError
)

type Logger struct {
	Formatter *LoggerFormatter
	Level     LoggerLevel
	Outs      []io.Writer
}

type LoggerFormatter struct{}

func New() *Logger {
	return &Logger{}
}

func Default() *Logger {
	logger := New()
	logger.Level = LevelDebug
	logger.Outs = append(logger.Outs, os.Stdout)
	logger.Formatter = &LoggerFormatter{}
	return logger
}

func (l *Logger) Error(msg string) {
	l.Print(LevelError, msg)
}

func (l *Logger) Info(msg string) {
	l.Print(LevelInfo, msg)
}

func (l *Logger) Debug(msg string) {
	l.Print(LevelDebug, msg)
}

func (l *Logger) Print(level LoggerLevel, msg string) {
	if l.Level > level {
		return
	}
	for _, out := range l.Outs {
		fmt.Fprint(out, msg)
	}
}

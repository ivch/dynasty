package logger

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// ErrParseStrToLevel indicates that string given to function ParseLevel can't be parsed to Level.
var ErrParseStrToLevel = errors.New("string can't be parsed to level, use: `error`, `info`, `warning`, `debug`")

// Level represents enumeration of logging levels.
type Level int

const (
	// DSB logger disabled.
	DSB Level = iota
	// ERR represents error logging level.
	ERR
	// WRN represents warning logging level.
	WRN
	// INF represents info logging level.
	INF
	// DBG represents debug logging level.
	DBG
)

func (l Level) String() string {
	levels := [5]string{"Disabled", "Error", "Warning", "Info", "Debug"}
	return levels[l]
}

// Logger represents logging logic.
type Logger interface {
	// Debug prints message with debug logging level.
	Debug(format string, v ...interface{})
	// Info prints message with info logging level.
	Info(format string, v ...interface{})
	// Warn prints message with warning logging level.
	Warn(format string, v ...interface{})
	// Error prints message with error logging level.
	Error(format string, v ...interface{})
}

// NewStdLog returns a new instance of StdLog struct.
// Takes variadic options which will be applied to StdLog.
func NewStdLog(opts ...Option) *StdLog {
	l := &StdLog{
		err: log.New(os.Stderr, "\033[31mERR\033[0m: ", log.Ldate|log.Ltime),
		wrn: log.New(os.Stderr, "\033[33mWRN\033[0m: ", log.Ldate|log.Ltime),
		inf: log.New(os.Stderr, "\033[32mINF\033[0m: ", log.Ldate|log.Ltime),
		dbg: log.New(os.Stderr, "\033[35mDBG\033[0m: ", log.Ldate|log.Ltime),
		lvl: INF,
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// StdLog represents standard library logger with levels.
type StdLog struct {
	err, wrn, inf, dbg *log.Logger
	lvl                Level
}

func (l *StdLog) Debug(format string, v ...interface{}) {
	if l.lvl < DBG {
		return
	}
	l.dbg.Printf(format, v...)
}

func (l *StdLog) Info(format string, v ...interface{}) {
	if l.lvl < INF {
		return
	}
	l.inf.Printf(format, v...)
}

func (l *StdLog) Warn(format string, v ...interface{}) {
	if l.lvl < WRN {
		return
	}
	l.wrn.Printf(format, v...)
}

func (l *StdLog) Error(format string, v ...interface{}) {
	if l.lvl < ERR {
		return
	}
	l.err.Printf(format, v...)
}

type Option func(l *StdLog)

func WithWriter(w io.Writer) Option {
	return func(l *StdLog) {
		l.err.SetOutput(w)
		l.wrn.SetOutput(w)
		l.inf.SetOutput(w)
		l.dbg.SetOutput(w)
	}
}

func WithLevel(level Level) Option {
	return func(l *StdLog) {
		l.lvl = level
	}
}

// ParseLevel parses string representation of logging level.
func ParseLevel(level string) (Level, error) {
	levels := map[string]Level{
		strings.ToLower(DSB.String()): DSB,
		strings.ToLower(ERR.String()): ERR,
		strings.ToLower(INF.String()): INF,
		strings.ToLower(WRN.String()): WRN,
		strings.ToLower(DBG.String()): DBG,
	}
	l, ok := levels[strings.ToLower(level)]
	if !ok {
		return INF, fmt.Errorf("%s %w", level, ErrParseStrToLevel)
	}
	return l, nil
}

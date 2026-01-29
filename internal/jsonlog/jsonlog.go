package jsonlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Level int8

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
	LevelOff
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

func (lg *Logger) print(level Level, msg string, properties map[string]string) (int, error) {
	if level < lg.minLevel {
		return 0, nil
	}

	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    msg,
		Properties: properties,
	}

	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	var line []byte

	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message: " + err.Error())
	}

	lg.mu.Lock()
	defer lg.mu.Unlock()

	return lg.out.Write(append(line, '\n'))
}

func (lg *Logger) Write(msg []byte) (int, error) {
	return lg.print(LevelError, string(msg), nil)
}

func (lg *Logger) PrintInfo(msg string, properties map[string]string) {
	lg.print(LevelInfo, msg, properties)
}

func (lg *Logger) PrintError(err error, properties map[string]string) {
	lg.print(LevelError, err.Error(), properties)
}

func (lg *Logger) PrintFatal(err error, properties map[string]string) {
	lg.print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

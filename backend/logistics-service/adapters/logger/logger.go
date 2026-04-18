package logger

import (
	"log/slog"
	"os"
	"strings"
)

type Logger struct {
	logger *slog.Logger
}

func New(logLevel string) *Logger {
	var level slog.Level
	switch strings.ToLower(logLevel) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	return &Logger{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})),
	}
}

func (l *Logger) Info(msg string, args ...any)  { l.logger.Info(msg, args...) }
func (l *Logger) Error(msg string, args ...any) { l.logger.Error(msg, args...) }
func (l *Logger) Debug(msg string, args ...any) { l.logger.Debug(msg, args...) }
func (l *Logger) Warn(msg string, args ...any)  { l.logger.Warn(msg, args...) }

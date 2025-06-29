package logger

import (
	"log/slog"
	"os"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
}

type logger struct {
	*slog.Logger
}

func New() Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	var handler slog.Handler
	format := os.Getenv("LOG_FORMAT")
	if format == "text" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return &logger{
		Logger: slog.New(handler),
	}
}

func (l *logger) Debug(msg string, args ...any) {
	l.Logger.Debug(msg, args...)
}

func (l *logger) Info(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.Logger.Warn(msg, args...)
}

func (l *logger) Error(msg string, args ...any) {
	l.Logger.Error(msg, args...)
}

func (l *logger) Fatal(msg string, args ...any) {
	l.Logger.Error(msg, args...)
	os.Exit(1)
}

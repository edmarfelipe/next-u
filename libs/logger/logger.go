package logger

import (
	"context"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type Logger interface {
	With(componentName string) Logger
	Info(ctx context.Context, msg string, keyvals ...interface{})
	Warn(ctx context.Context, msg string, keyvals ...interface{})
	Error(ctx context.Context, msg string, keyvals ...interface{})
}

type Field struct {
	Key   string
	Value any
}

func New() Logger {
	w := log.NewSyncWriter(os.Stderr)

	return &logger{
		instance: log.NewJSONLogger(w),
	}
}

type logger struct {
	instance log.Logger
}

func (l *logger) With(componentName string) Logger {
	return &logger{
		instance: log.With(l.instance, "component", componentName),
	}
}

func (l *logger) Info(ctx context.Context, msg string, keyvals ...interface{}) {
	level.Info(l.instance).Log("msg", msg, keyvals)
}

func (l *logger) Warn(ctx context.Context, msg string, keyvals ...interface{}) {
	level.Warn(l.instance).Log("msg", msg, keyvals)
}

func (l *logger) Error(ctx context.Context, msg string, keyvals ...interface{}) {
	level.Error(l.instance).Log("msg", msg, keyvals)
}

package logger

import (
	"context"
	"os"

	"github.com/go-kit/log"
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
	lg := log.NewJSONLogger(log.NewSyncWriter(os.Stderr))

	return &logger{
		instance: log.WithPrefix(lg, "ts", log.DefaultTimestampUTC, "caller", log.Caller(4)),
	}
}

type logger struct {
	instance log.Logger
}

func (l *logger) With(componentName string) Logger {
	return &logger{
		instance: log.WithPrefix(l.instance, "component", componentName),
	}
}

func (l *logger) Info(ctx context.Context, msg string, keyvals ...interface{}) {
	log.WithPrefix(l.instance, "level", "info", "rqid", getRequestID(ctx), "msg", msg).Log(keyvals...)
}

func (l *logger) Warn(ctx context.Context, msg string, keyvals ...interface{}) {
	log.WithPrefix(l.instance, "level", "info", "rqid", getRequestID(ctx), "msg", msg).Log(keyvals...)
}

func (l *logger) Error(ctx context.Context, msg string, keyvals ...interface{}) {
	log.WithPrefix(l.instance, "level", "error", "rqid", getRequestID(ctx), "msg", msg).Log(keyvals...)
}

func getRequestID(ctx context.Context) string {
	value := ctx.Value("requestid")
	if id, ok := value.(string); ok {
		return id
	}
	return "unknown"
}

package log

import "context"

const (
	CorrelationID = "correlation_id"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	With(args ...interface{}) Logger
	WithCtx(ctx context.Context, args ...interface{}) Logger
	Infow(string, ...interface{})
	Named(string) Logger
}

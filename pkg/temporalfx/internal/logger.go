package temporal

import (
	"fmt"
	"github.com/postmanq/postmanq/pkg/logfx/log"

	"go.uber.org/zap"
)

type logAdapter struct {
	logger log.Logger
}

func (l *logAdapter) fields(keyvals []interface{}) []zap.Field {
	if len(keyvals)%2 != 0 {
		return []zap.Field{zap.Error(fmt.Errorf("odd number of keyvals pairs: %v", keyvals))}
	}

	var fields []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			key = fmt.Sprintf("%v", keyvals[i])
		}
		fields = append(fields, zap.Any(key, keyvals[i+1]))
	}

	return fields
}

func (l *logAdapter) Debug(msg string, keyvals ...interface{}) {
	l.logger.Debug(msg, l.fields(keyvals))
}

func (l *logAdapter) Info(msg string, keyvals ...interface{}) {
	l.logger.Info(msg, l.fields(keyvals))
}

func (l *logAdapter) Warn(msg string, keyvals ...interface{}) {
	l.logger.Warn(msg, l.fields(keyvals))
}

func (l *logAdapter) Error(msg string, keyvals ...interface{}) {
	l.logger.Error(msg, l.fields(keyvals))
}

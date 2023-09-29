package services

import (
	"bufio"
	"context"
	"github.com/postmanq/postmanq/pkg/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

var (
	keys = []string{
		log.CorrelationID,
	}
)

func NewDefaultConfig(lvl ...string) zap.Config {
	level := zap.NewAtomicLevel()
	err := level.UnmarshalText([]byte(strings.Join(lvl, "")))
	if err != nil || len(lvl) == 0 {
		level.SetLevel(zap.DebugLevel)
	}
	return zap.Config{
		Level:            level,
		DisableCaller:    false,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		Encoding:         "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			NameKey:        "logger",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			CallerKey:      "file",
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
			TimeKey:        "time",
		},
	}
}

func NewTestLogger(writer *bufio.Writer) log.Logger {
	cfg := NewDefaultConfig()
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	l := zap.New(
		zapcore.NewCore(encoder, zapcore.AddSync(writer), zapcore.DebugLevel)).
		Sugar()

	return &logger{
		Config:        &cfg,
		SugaredLogger: l,
	}
}

func NewLogger(
	cfg zap.Config,
	fields ...zap.Field,
) (log.Logger, error) {
	baseLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	l := baseLogger.WithOptions(
		zap.Fields(fields...),
	).Sugar()

	defer l.Sync() //nolint:errcheck

	return &logger{
		Config:        &cfg,
		SugaredLogger: l,
	}, nil
}

type logger struct {
	*zap.Config
	*zap.SugaredLogger
}

func (l *logger) With(args ...interface{}) log.Logger {
	return &logger{
		Config:        l.Config,
		SugaredLogger: l.SugaredLogger.With(args...),
	}
}

func (l *logger) WithCtx(ctx context.Context, args ...interface{}) log.Logger {
	logger := &logger{
		Config:        l.Config,
		SugaredLogger: l.SugaredLogger.With(args...),
	}

	for _, key := range keys {
		value, ok := ctx.Value(key).(interface{})
		if ok {
			args = append(args, key, value)
		}
	}

	return logger.With(args...)
}

func (l *logger) Named(name string) log.Logger {
	return &logger{
		Config:        l.Config,
		SugaredLogger: l.SugaredLogger.Named(name),
	}
}

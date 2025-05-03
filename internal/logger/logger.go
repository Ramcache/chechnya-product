package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func NewLogger() (*zap.Logger, error) {
	// INFO + WARN → logs/info.2025-05-03.log
	infoWriter, err := rotatelogs.New(
		"logs/info.%Y-%m-%d.log",
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		return nil, err
	}

	// ERROR → logs/error.2025-05-03.log
	errorWriter, err := rotatelogs.New(
		"logs/error.%Y-%m-%d.log",
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		return nil, err
	}

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		MessageKey:   "message",
		CallerKey:    "caller",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	infoLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zap.InfoLevel && l < zap.ErrorLevel
	})

	errorLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zap.ErrorLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.AddSync(errorWriter), errorLevel),
	)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)), nil
}

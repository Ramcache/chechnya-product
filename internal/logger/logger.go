package logger

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() (*zap.Logger, error) {
	// 📂 Файлы для info, error, debug
	infoWriter, err := rotatelogs.New(
		"logs/info.%Y-%m-%d.log",
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		return nil, err
	}

	errorWriter, err := rotatelogs.New(
		"logs/error.%Y-%m-%d.log",
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		return nil, err
	}

	debugWriter, err := rotatelogs.New(
		"logs/debug.%Y-%m-%d.log",
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		return nil, err
	}

	// ⚙️ Настройка форматирования
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		MessageKey:   "message",
		CallerKey:    "caller",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	// 🔎 Уровни логирования
	infoLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zap.InfoLevel && l < zap.ErrorLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zap.ErrorLevel
	})
	debugLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l == zap.DebugLevel
	})

	// 🧩 Комбинируем все уровни
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.AddSync(errorWriter), errorLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.AddSync(debugWriter), debugLevel),
	)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)), nil
}

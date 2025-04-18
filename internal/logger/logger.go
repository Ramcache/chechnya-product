package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// NewLogger создаёт новый zap.Logger с выводом в stdout и файл
func NewLogger() (*zap.Logger, error) {
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, err
	}

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:         "json",
		OutputPaths:      []string{"stdout", "logs/app.log"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:      "time",
			LevelKey:     "level",
			MessageKey:   "message",
			CallerKey:    "caller",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	return cfg.Build()
}

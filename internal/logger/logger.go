// internal/logger/logger.go
package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func NewLogger() (*zap.Logger, error) {
	writer, err := rotatelogs.New(
		"logs/app.%Y-%m-%d.log",                   // шаблон имени
		rotatelogs.WithLinkName("logs/app.log"),   // симлинк на последний файл
		rotatelogs.WithMaxAge(7*24*time.Hour),     // хранить 7 дней
		rotatelogs.WithRotationTime(24*time.Hour), // ротация каждый день
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

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(writer),
		zap.InfoLevel,
	)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)), nil
}

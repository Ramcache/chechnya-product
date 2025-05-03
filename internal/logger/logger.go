package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func NewLogger() (*zap.Logger, error) {
	infoWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/info.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
	})

	errorWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/error.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
	})

	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		MessageKey:   "message",
		CallerKey:    "caller",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	})

	infoLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zap.InfoLevel && l < zap.ErrorLevel
	})

	errorLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zap.ErrorLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoWriter, zapcore.AddSync(os.Stdout)), infoLevel),
		zapcore.NewCore(encoder, errorWriter, errorLevel),
	)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)), nil
}

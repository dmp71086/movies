package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Global *zap.Logger

func InitLogger() {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	level := zap.NewAtomicLevelAt(zapcore.WarnLevel)
	if os.Getenv("DEBUG") == "true" {
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	// cfg := zap.Config{
	// 	Encoding:         "json",
	// 	Level:            zap.NewAtomicLevelAt(zapcore.ErrorLevel),
	// 	OutputPaths:      []string{"stdoyt"},
	// 	ErrorOutputPaths: []string{"stderr"},
	// 	EncoderConfig:    encoderConfig,
	// }

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // Custom encoder
		zapcore.AddSync(os.Stdout),            // Sync to stdout
		level,
	)

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	// core := zapcore.NewTee(
	// 	zapcore.NewCore(kafkaEncoder, topicErrors, highPriority),
	// 	zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
	// 	zapcore.NewCore(kafkaEncoder, topicDebugging, lowPriority),
	// 	zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	// )

	log := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.FatalLevel))

	log.Info("This is an INFO message.Logger works in debug mode")

	Global = log
}

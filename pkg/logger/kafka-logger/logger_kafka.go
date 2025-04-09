package kafka_logger

import (
	"io"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type KafkaWriter struct {
	topic string
}

// NewZapLogger create new zap logger with custom config
func NewZapLogger(iw io.Writer) *zap.Logger {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	topicErrors := zapcore.AddSync(iw)
	kafkaEncoder := zap.NewProductionEncoderConfig()
	kafkaEncoder.TimeKey = "timestamp"
	kafkaEncoder.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(kafkaEncoder), topicErrors, highPriority),
	)

	return zap.New(core, zap.AddCaller())
}

// NewKafkaWriter create new io.Writer implementation for custom zap logger
func NewKafkaWriter(topic string) *KafkaWriter {
	return &KafkaWriter{
		topic: topic,
	}
}

// Write bytes to kafka
func (k *KafkaWriter) Write(p []byte) (n int, err error) {
	n = 0

	err = Produce(k.topic, p)
	if err != nil {
		log.Fatal("failed to write messages: ", err.Error())
	}

	return
}

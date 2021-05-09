package zapflatencoder

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	err := zap.RegisterEncoder(EncoderName, func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return newEncoder(cfg), nil
	})
	if err != nil {
		panic(err)
	}
}

// DefaultEncoder returns default console encoder
func DefaultEncoder() zapcore.Encoder {
	return newEncoder(DefaultConfig)
}

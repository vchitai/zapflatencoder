package main

import (
	"fmt"
	"log"

	"github.com/vchitai/zapflatencoder"
	"go.uber.org/zap"
)

func main() {
	l, err := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      false,
		Encoding:         zapflatencoder.EncoderName,
		EncoderConfig:    zapflatencoder.DefaultConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()
	if err != nil {
		log.Fatal("Cannot init logger", err)
	}
	l.Info("Hello World")
	l.Debug("Hello World")
	l.Error("Hello World with error", zap.Error(fmt.Errorf("this is an error")))
	l.Warn("Hello World with error", zap.Error(fmt.Errorf("this is an error")))
	l.DPanic("Hello World with error", zap.Error(fmt.Errorf("this is an error")))
}

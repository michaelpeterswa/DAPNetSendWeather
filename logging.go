package main

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	logger.Info("DAPNetSendWeather is initializing...")
}

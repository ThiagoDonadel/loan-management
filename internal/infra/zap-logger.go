package infra

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func InitlializeLogger() {

	if viper.Get("profile") == "local" {
		logger, _ := zap.NewDevelopment()
		Logger = logger.Sugar()
	} else {
		logger, _ := zap.NewProduction()
		Logger = logger.Sugar()
	}

}

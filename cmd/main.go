package main

import (
	"telebotai/pkg/service"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	if err := initConfig(); err != nil {
		logrus.Fatalf("cannot parse configs, %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("can't load environment, %s", err.Error())
	}

	services := service.Service{}
	services.Run()
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

package main

import (
	"gw-notification/iternal/broker/kafka"
	"gw-notification/iternal/config"
	mongodb "gw-notification/iternal/storage/mongo"
	"gw-notification/pkg/logging"
	"os"
)

func main() {

	//0. Подготовка логгера
	logger := logging.GetLogger()
	logger.Info("Start...")

	//1. Загрузка конфига
	cfg := config.GetConfig()
	logger.Infof("Get Config: %v", cfg)

	//2. Подключение к mongoDB

	mongoStorage, err := mongodb.NewStorage(cfg.MongoConn)
	logger.Infof("Подключение к mongoDB: %s", cfg.MongoConn)
	if err != nil {
		logger.Error("Failed to connect to MongoDB", "error", err)
		os.Exit(1)
	}
	defer mongoStorage.Close()

	if err = kafka.StartConsumer(cfg, logger, mongoStorage); err != nil {
		logger.Error("Consumer failed", "error", err)
	}
}

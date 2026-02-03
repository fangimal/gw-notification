package config

import (
	"gw-notification/pkg/logging"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	MongoConn string `yaml:"storage" env:"storage" env-default:"mongodb://localhost:27018"`

	Kafka KafkaConfig `yaml:"kafka"`
}

type KafkaConfig struct {
	KafkaServerAddress string `yaml:"KafkaServerAddress" env:"KafkaServerAddress" env-default:"localhost:9092"`
	KafkaTopic         string `yaml:"KafkaTopic" env:"KafkaTopic" env-default:"notification"`
	KafkaGroupID       string `yaml:"KafkaGroupID" env:"KafkaGroupID" env-default:"notifications"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	logger := logging.GetLogger()

	once.Do(func() {
		instance = &Config{}

		// Сначала пытаемся прочитать конфигурацию из YAML-файла
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			logger.Warnf("Could not read config file, using defaults and environment variables: %v", err)
		}

		// Затем загружаем переменные окружения, которые переопределят значения из YAML
		if err := cleanenv.ReadEnv(instance); err != nil {
			logger.Warnf("Warning: Error reading config from environment: %v", err)
		}
	})
	return instance
}

package config

import (
	"gw-notification/pkg/logging"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	MongoConn string `yaml:"mongo_conn" env-default:"mongodb://localhost:27018"`

	Kafka KafkaConfig `yaml:"kafka"`
}

type KafkaConfig struct {
	KafkaServerAddress string `yaml:"KafkaServerAddress" env-default:"localhost:9092"`
	KafkaTopic         string `yaml:"KafkaTopic" env-default:"notification"`
	KafkaGroupID       string `yaml:"KafkaGroupID" env-default:"notifications"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	logger := logging.GetLogger()

	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Infof("Error reading config: %v", help)
		}
	})
	return instance
}

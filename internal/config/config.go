package config

import (
	"rest/internal/logging"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug *bool `yaml:is_debug`
	Listen  struct {
		Type string `yaml:type`
		Port string `yaml:port`
	} `yaml:listen`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logging.GetLogger().Info("read application cofigurations")
		instance := &Config{}
		logging.GetLogger().Info("go to file")
		if err := cleanenv.ReadConfig("/Users/samokat/learn/rest_learn/config.yaml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}

	})
	return instance
}

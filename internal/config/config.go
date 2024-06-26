package config

import (
	"rest/internal/logging"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug" env-required:"true"`
	Listen  struct {
		Type     string `yaml:"type" env-default:"port"`
		BindIP   string `yaml:"bind_ip" env-default:"127.0.0.1"`
		HttpPort string `yaml:"http_port" env-default:"8080"`
		GrpcPort string `yaml:"grpc_port" env-default:"8080"`
		URI_List string `yaml:"URI_List" env-default:"/"`
		URI_Once string `yaml:"URI_Once" env-default:"/"`
		HOST     string `yaml:"HOST" env-default:"/"`
	} `yaml:"listen"`
	Storage StorageConfig `yaml:"storage"`
	User    UserConfig    `yaml:"User"`
}

type StorageConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserConfig struct {
	Name  string `yaml:"Name" env-default:"/"`
	Job   string `yaml:"Job" env-default:"/"`
	MinId int    `yaml:"MinId" env-default:"/"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("/Users/samokat/learn/rest_learn/config.yaml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}

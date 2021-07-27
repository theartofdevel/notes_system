package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/theartofdevel/notes_system/file_service/pkg/logging"
	"sync"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"localhost"`
		Port   string `yaml:"port" env-default:"10002"`
	}
	MinIO struct {
		Endpoint  string `yaml:"endpoint" env-required:"true"`
		AccessKey string `yaml:"access_key" env-required:"true"`
		SecretKey string `yaml:"secret_key" env-required:"true"`
	} `yaml:"minio"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}

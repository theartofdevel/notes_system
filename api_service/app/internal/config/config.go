package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/theartofdevel/notes_system/api_service/pkg/logging"
	"sync"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	JWT     struct {
		Secret string `yaml:"secret" env-required:"true"`
	}
	Listen struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"localhost"`
		Port   string `yaml:"port" env-default:"8080"`
	}
	CategoryService struct {
		URL string `yaml:"url" env-required:"true"`
	} `yaml:"category_service" env-required:"true"`
	UserService struct {
		URL string `yaml:"url" env-required:"true"`
	} `yaml:"user_service" env-required:"true"`
	NoteService struct {
		URL string `yaml:"url" env-required:"true"`
	} `yaml:"note_service" env-required:"true"`
	TagService struct {
		URL string `yaml:"url" env-required:"true"`
	} `yaml:"tag_service" env-required:"true"`
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

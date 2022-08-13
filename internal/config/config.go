package config

import (
	"auth_pd/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	IsDebug bool   `yaml:"is_debug"`
	LogsDir string `yaml:"logs_dir"`
	Listen  struct {
		BindIp string `yaml:"bind_ip"`
		Port   string `yaml:"port"`
	} `yaml:"listen"`
	Database struct {
		Scheme   string `yaml:"scheme"`
		Login    string `yaml:"login"`
		Password string `yaml:"password"`
		BindIp   string `yaml:"bind_ip"`
		Port     string `yaml:"port"`
	} `yaml:"database"`
}

var cfg *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Infof("Чтение конфига")
		cfg = &Config{}
		err := cleanenv.ReadConfig("config.yaml", cfg)
		if err != nil {
			logger.Fatalf("Не удалось вычитать конфиг: %s", err.Error())
		}
	})
	return cfg
}

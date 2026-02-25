package modules

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type PostgreConfig struct {
	Host        string        `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port        string        `yaml:"port" env:"DB_PORT" env-default:"5434"`
	Username    string        `yaml:"username" env:"DB_USER" env-default:"postgres"`
	Password    string        `yaml:"password" env:"DB_PASSWORD" env-default:"postgres"`
	DBName      string        `yaml:"db_name" env:"DB_NAME" env-default:"go_kbtu"`
	SSLMode     string        `yaml:"ssl_mode" env:"DB_SSLMODE" env-default:"disable"`
	ExecTimeout time.Duration `yaml:"exec_timeout" env:"DB_EXEC_TIMEOUT" env-default:"5s"`
}

type AppConfig struct {
	Port   string        `yaml:"app_port" env:"APP_PORT" env-default:"8080"`
	APIKey string        `yaml:"api_key" env:"API_KEY" env-default:"secret12345"`
	PG     PostgreConfig `yaml:"db"`
}

func LoadConfig(path string) (*AppConfig, error) {
	var cfg AppConfig
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}
	return &cfg, nil
}

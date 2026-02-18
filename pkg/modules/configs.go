package modules

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type PostgreConfig struct {
	Host        string
	Port        string
	Username    string
	Password    string
	DBName      string
	SSLMode     string
	ExecTimeout time.Duration
}

type AppConfig struct {
	Port   string
	APIKey string
	PG     *PostgreConfig
}

func LoadConfig() *AppConfig {
	_ = godotenv.Load()

	timeoutSec := getIntEnv("DB_EXEC_TIMEOUT_SEC", 5)

	return &AppConfig{
		Port:   getEnv("APP_PORT", "8080"),
		APIKey: getEnv("API_KEY", "secret12345"),
		PG: &PostgreConfig{
			Host:        getEnv("DB_HOST", "localhost"),
			Port:        getEnv("DB_PORT", "5434"),
			Username:    getEnv("DB_USER", "postgres"),
			Password:    getEnv("DB_PASSWORD", "postgres"),
			DBName:      getEnv("DB_NAME", "go_kbtu"),
			SSLMode:     getEnv("DB_SSLMODE", "disable"),
			ExecTimeout: time.Duration(timeoutSec) * time.Second,
		},
	}
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func getIntEnv(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

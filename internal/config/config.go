package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"` // Если параметр не указан, то по умолчанию используется "local"
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-default:"1h"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// Функция не будет возращать ошибку, а будет паниковать в случае ошибки
func MustLoad() *Config {
	err := godotenv.Load() // Работаем с .env файлом

	if err != nil {
		panic("Error loading .env file")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config " + err.Error())
	}

	return &cfg
}

package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server      Server
	DB          DB          `json:"database"`
	ExternalSvc ExternalSvc `json:"external_service"`
	EnableRedis bool        `env:"ENABLE_REDIS"`
	RedisURL    string      `env:"REDIS_URL"`
}

type Server struct {
	Port string
}

type DB struct {
	URL     string `env:"PG_URL"`
	Host    string
	Port    string
	User    string
	Pass    string
	Name    string
	PoolMax int    `json:"pool_max"`
	SslMode string `json:"ssl_mode"`
}

type ExternalSvc struct {
	Url string
	N   int
	K   int
	X   int
}

func LoadConfig(file string) (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadConfig(file, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

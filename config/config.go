package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server         Server
	DB             DB      `json:"database"`
	ExternalSvcUrl string  `json:"external_svc_url"`
	TaskOps        TaskOps `json:"task_ops"`
	EnableRedis    bool    `env:"ENABLE_REDIS"`
	Redis          Redis
}

type Redis struct {
	Host string `env:"REDIS_HOST" env-default:"localhost"`
	Port string `env:"REDIS_PORT" env-default:"6379"`
	Pass string `env:"REDIS_PASSWORD"`
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

type TaskOps struct {
	// N is count of requests per K seconds
	N int
	// K is time established for N requests
	K int
	// X blocking time when limits are exceeded
	X int
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

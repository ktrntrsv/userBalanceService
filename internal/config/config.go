package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server   `yaml:"server"`
	Postgres `yaml:"postgres"`
	Logger   `yaml:"logger"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password" env:"POSTGRES_PWD"`
	Database string `yaml:"database"`
}

type Server struct {
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

type Logger struct {
	Level string `yaml:"level"`
}

func NewConfig(confPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(confPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

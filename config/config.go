package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App `yaml:"app"`
		TCP `yaml:"http"`
		Log `yaml:"logger"`
		AoF `yaml:"aof"`
	}

	TCP struct {
		Addr string `env-required:"true" yaml:"addr" env:"ADDR"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level" env:"LOG_LEVEL"`
	}

	AoF struct {
		FileName string `env-required:"true" yaml:"file_name" env:"FILE_NAME"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	config := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", config)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(config)
	if err != nil {
		return nil, fmt.Errorf("Env read error: %w", err)
	}

	return config, nil
}

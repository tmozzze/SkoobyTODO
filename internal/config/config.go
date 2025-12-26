package config

import (
	"fmt"
	"os"

	"github.com/tmozzze/SkoobyTODO/internal/utils"
)

type Config struct {
	Env string
}

func New() *Config {
	return &Config{}
}

func (c *Config) Load(envPath string) error {
	const op = "config.config.Load"
	err := utils.LoadEnv(envPath)
	if err != nil {
		return fmt.Errorf("%s: load config failed: %w", op, err)
	}

	c.Env = os.Getenv("ENV")
	return nil
}

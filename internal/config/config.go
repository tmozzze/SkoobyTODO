package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/tmozzze/SkoobyTODO/internal/utils"
)

type Config struct {
	Env string

	// server
	ServerPort string
	// timeouts
	ReadTimeout    int
	WriteTimeout   int
	IdleTimeout    int
	HandlerTimeout int
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
	c.ServerPort = os.Getenv("SERVER_PORT")

	c.ReadTimeout, err = strconv.Atoi(os.Getenv("READ_TIMEOUT"))
	if err != nil {
		return fmt.Errorf("%s: ReadTimeout conv to int error: %w", op, err)
	}

	c.WriteTimeout, err = strconv.Atoi(os.Getenv("WRITE_TIMEOUT"))
	if err != nil {
		return fmt.Errorf("%s: WriteTimeout conv to int error: %w", op, err)
	}

	c.IdleTimeout, err = strconv.Atoi(os.Getenv("IDLE_TIMEOUT"))
	if err != nil {
		return fmt.Errorf("%s: IdleTimeout conv to int error: %w", op, err)
	}

	c.HandlerTimeout, err = strconv.Atoi(os.Getenv("HANDLER_TIMEOUT"))
	if err != nil {
		return fmt.Errorf("%s: HandlerTimeout conv to int error: %w", op, err)
	}

	fmt.Println("Config loaded successfully")

	return nil
}

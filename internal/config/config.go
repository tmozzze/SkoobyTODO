package config

import (
	"log"
	"os"

	"github.com/tmozzze/SkoobyTODO/internal/utils"
)

type Config struct {
	VarEnv string
}

func New() *Config {
	err := utils.LoadEnv(".env")
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	varEnv := os.Getenv("MY_VAR")
	return &Config{VarEnv: varEnv}
}

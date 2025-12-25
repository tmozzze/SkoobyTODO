package main

import (
	"fmt"

	"github.com/tmozzze/SkoobyTODO/internal/config"
)

func main() {
	fmt.Println("Hello, world!")

	cfg := config.New()

	fmt.Println(cfg.VarEnv)
}

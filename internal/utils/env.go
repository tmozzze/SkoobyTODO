package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func LoadEnv(path string) error {
	const op = "utils.env.LoadEnv"

	if path == "" {
		return fmt.Errorf("%s: env open failed: %w", op, errors.New("path is empty"))
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("%s: env open failed: %w", op, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// split key, value and set env
		parts := strings.SplitN(line, "=", 2)

		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			os.Setenv(key, val)
		}
	}

	// check scanner error

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("%s: scanner error: %w", op, err)
	}

	return nil

}

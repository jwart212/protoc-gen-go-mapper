package env

import (
	"errors"
	"fmt"
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Options struct {
	Prefix string

	DotEnv bool

	DotEnvFiles []string
}

func MustLoad[T any](
	opt Options,
) (*T, error) {

	cfg := new(T)

	if opt.DotEnv {

		files := opt.DotEnvFiles

		if len(files) == 0 {
			files = []string{
				".env",
				".env.local",
			}
		}

		_ = loadDotEnv(files...)
	}

	if err := env.ParseWithOptions(
		cfg,
		env.Options{
			Prefix: opt.Prefix,
		},
	); err != nil {
		return nil, fmt.Errorf(
			"parse environment: %w",
			err,
		)
	}

	return cfg, nil
}

func loadDotEnv(files ...string) error {

	for _, file := range files {

		if _, err := os.Stat(file); err != nil {

			if errors.Is(err, os.ErrNotExist) {
				continue
			}

			return err
		}

		if err := godotenv.Overload(file); err != nil {
			return err
		}
	}

	return nil
}

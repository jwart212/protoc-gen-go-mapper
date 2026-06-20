package config_test

import (
	"fmt"

	"github.com/jwart212/protoc-gen-go-mapper/internal/config"
)

func ExampleConfig() {
	cfg := &config.Config{
		Version:  "v1",
		Database: "sqlc",
		Package: config.Package{
			Proto: "internal/gen",
			DB:    "internal/postgres",
		},
	}
	fmt.Printf("Database: %s, Proto Package: %s", cfg.Database, cfg.Package.Proto)
	// Output: Database: sqlc, Proto Package: internal/gen
}

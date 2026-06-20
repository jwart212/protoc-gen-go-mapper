package plugin_test

import (
	"fmt"

	"github.com/jwart212/protoc-gen-go-mapper/internal/config"
	"github.com/jwart212/protoc-gen-go-mapper/internal/plugin"
)

func ExampleNew() {
	cfg := &config.Config{
		Version:  "v1",
		Database: "sqlc",
		Package: config.Package{
			Proto: "internal/gen",
			DB:    "internal/postgres",
		},
	}

	p := plugin.New(cfg)
	fmt.Printf("Plugin created: %v", p != nil)
	// Output: Plugin created: true
}

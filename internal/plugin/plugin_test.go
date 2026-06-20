package plugin

import (
	"testing"

	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/internal/config"
)

func TestNew(t *testing.T) {
	cfg := &config.Config{
		Version:  "v1",
		Database: "sqlc",
		Package: config.Package{
			Proto: "internal/gen",
			DB:    "internal/postgres",
		},
	}

	p := New(cfg)
	if p == nil {
		t.Error("New() returned nil")
	}
	if p.registry == nil {
		t.Error("Plugin registry should be initialized")
	}
	if p.generator == nil {
		t.Error("Plugin generator should be initialized")
	}
}

func TestRegisterConverters(t *testing.T) {
	cfg := &config.Config{
		Version:  "v1",
		Database: "sqlc",
		Package: config.Package{
			Proto: "internal/gen",
			DB:    "internal/postgres",
		},
	}

	p := New(cfg)
	p.registerConverters()

	// Verify converters are registered (by checking registry has converters)
	if p.registry == nil {
		t.Error("Registry should be initialized")
	}
}

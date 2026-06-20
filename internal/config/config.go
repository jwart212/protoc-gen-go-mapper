package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/jwart212/protoc-gen-go-mapper/pkg/errors"
)

// Config represents the mapper.yaml configuration.
type Config struct {
	Version      string            `yaml:"version"`
	Database     string            `yaml:"database"`
	DBPackage    string            `yaml:"db_package"` // Direct DB package import path
	Package      Package           `yaml:"package"`
	TypeMappings map[string]string `yaml:"type_mappings"` // Proto message name to DB type name mapping
	Messages     []string          `yaml:"messages"`      // List of message names to generate mappers for
}

// Package represents the package configuration.
type Package struct {
	Proto string `yaml:"proto"`
	DB    string `yaml:"db"`
}

// Load reads and parses the mapper.yaml configuration from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %s: %w", path, err)
	}

	if err := Validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate checks that the configuration is valid.
func Validate(cfg *Config) error {
	if cfg.Version == "" {
		return fmt.Errorf("validating config: version: %w", errors.ErrInvalidConfig)
	}

	if cfg.Version != "v1" {
		return fmt.Errorf("validating config: version %s: %w", cfg.Version, errors.ErrInvalidConfig)
	}

	validDatabases := map[string]bool{
		"sqlc":         true,
		"pgx":          true,
		"database_sql": true,
	}

	if cfg.Database == "" {
		return fmt.Errorf("validating config: database: %w", errors.ErrInvalidConfig)
	}

	if !validDatabases[cfg.Database] {
		return fmt.Errorf("validating config: database %s: %w", cfg.Database, errors.ErrInvalidConfig)
	}

	if cfg.Package.Proto == "" {
		return fmt.Errorf("validating config: package.proto: %w", errors.ErrInvalidConfig)
	}

	if cfg.Package.DB == "" {
		return fmt.Errorf("validating config: package.db: %w", errors.ErrInvalidConfig)
	}

	return nil
}

package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	mappererrors "github.com/jwart212/protoc-gen-go-mapper/pkg/errors"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &Config{
				Version:  "v1",
				Database: "sqlc",
				Package: Package{
					Proto: "internal/gen",
					DB:    "internal/postgres",
				},
			},
			wantErr: false,
		},
		{
			name: "missing version",
			cfg: &Config{
				Database: "sqlc",
				Package: Package{
					Proto: "internal/gen",
					DB:    "internal/postgres",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid version",
			cfg: &Config{
				Version:  "v2",
				Database: "sqlc",
				Package: Package{
					Proto: "internal/gen",
					DB:    "internal/postgres",
				},
			},
			wantErr: true,
		},
		{
			name: "missing database",
			cfg: &Config{
				Version: "v1",
				Package: Package{
					Proto: "internal/gen",
					DB:    "internal/postgres",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid database",
			cfg: &Config{
				Version:  "v1",
				Database: "invalid",
				Package: Package{
					Proto: "internal/gen",
					DB:    "internal/postgres",
				},
			},
			wantErr: true,
		},
		{
			name: "missing package.proto",
			cfg: &Config{
				Version:  "v1",
				Database: "sqlc",
				Package: Package{
					DB: "internal/postgres",
				},
			},
			wantErr: true,
		},
		{
			name: "missing package.db",
			cfg: &Config{
				Version:  "v1",
				Database: "sqlc",
				Package: Package{
					Proto: "internal/gen",
				},
			},
			wantErr: true,
		},
		{
			name: "valid pgx config",
			cfg: &Config{
				Version:  "v1",
				Database: "pgx",
				Package: Package{
					Proto: "internal/gen",
					DB:    "internal/postgres",
				},
			},
			wantErr: false,
		},
		{
			name: "valid database_sql config",
			cfg: &Config{
				Version:  "v1",
				Database: "database_sql",
				Package: Package{
					Proto: "internal/gen",
					DB:    "internal/postgres",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, mappererrors.ErrInvalidConfig) {
				t.Errorf("Validate() error should wrap ErrInvalidConfig, got %v", err)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("valid config file", func(t *testing.T) {
		configPath := filepath.Join(tmpDir, "mapper.yaml")
		configContent := `version: v1
database: sqlc
package:
  proto: internal/gen
  db: internal/postgres
`
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatal(err)
		}

		cfg, err := Load(configPath)
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		if cfg.Version != "v1" {
			t.Errorf("Expected version v1, got %s", cfg.Version)
		}
		if cfg.Database != "sqlc" {
			t.Errorf("Expected database sqlc, got %s", cfg.Database)
		}
		if cfg.Package.Proto != "internal/gen" {
			t.Errorf("Expected package.proto internal/gen, got %s", cfg.Package.Proto)
		}
		if cfg.Package.DB != "internal/postgres" {
			t.Errorf("Expected package.db internal/postgres, got %s", cfg.Package.DB)
		}
	})

	t.Run("invalid config file", func(t *testing.T) {
		configPath := filepath.Join(tmpDir, "invalid.yaml")
		configContent := `version: v2
database: sqlc
package:
  proto: internal/gen
  db: internal/postgres
`
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatal(err)
		}

		_, err := Load(configPath)
		if err == nil {
			t.Error("Load() should return error for invalid config")
		}
		if !errors.Is(err, mappererrors.ErrInvalidConfig) {
			t.Errorf("Load() error should wrap ErrInvalidConfig, got %v", err)
		}
	})

	t.Run("file not found", func(t *testing.T) {
		configPath := filepath.Join(tmpDir, "nonexistent.yaml")
		_, err := Load(configPath)
		if err == nil {
			t.Error("Load() should return error for nonexistent file")
		}
	})
}

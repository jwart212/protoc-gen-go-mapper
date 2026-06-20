package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/jwart212/protoc-gen-go-mapper/pkg/errors"
)

// Config represents the mapper.yaml configuration.
type Config struct {
	Version              string                     `yaml:"version"`
	Database             string                     `yaml:"database"`
	DBPackage            string                     `yaml:"db_package"` // Direct DB package import path
	Package              Package                    `yaml:"package"`
	TypeMappings         map[string]string          `yaml:"type_mappings"`          // Proto message name to DB type name mapping
	ResponseTypeMappings map[string]string          `yaml:"response_type_mappings"` // Response message to SQLC Row type mapping
	Messages             []string                   `yaml:"messages"`               // List of message names to generate mappers for
	FieldHandlers        []FieldHandlerConfig       `yaml:"field_handlers"`         // Field-level handler configurations
	TypeConversions      []TypeConversionConfig     `yaml:"type_conversions"`       // Type-based conversion configurations
	HelperFunctions      []HelperFunctionConfig     `yaml:"helper_functions"`       // Custom helper function definitions
	ResponsePatterns     ResponsePatternsConfig     `yaml:"response_patterns"`      // Response field pattern configuration
	TypeAliases          map[string]TypeAliasConfig `yaml:"type_aliases"`           // Reusable type conversion aliases
	PointerSettings      PointerSettingsConfig      `yaml:"pointer_settings"`       // Pointer handling strategies
}

// FieldHandlerConfig represents configuration for a field handler.
type FieldHandlerConfig struct {
	Name         string   `yaml:"name"`           // Handler name for identification
	Type         string   `yaml:"type"`           // Handler type: skip, type_assertion, default_value, field_mapping
	MatchField   string   `yaml:"match_field"`    // Field name to match (case-insensitive)
	MatchDBTypes []string `yaml:"match_db_types"` // DB type names to match
	MatchMessage string   `yaml:"match_message"`  // Message name to match (case-insensitive)
	AssertType   string   `yaml:"assert_type"`    // Type to assert to (for type_assertion)
	DefaultValue string   `yaml:"default_value"`  // Default value expression (for default_value)
	ToProto      string   `yaml:"to_proto"`       // Custom expression for DB -> Proto (for field_mapping)
	ToDB         string   `yaml:"to_db"`          // Custom expression for Proto -> DB (for field_mapping)
	Priority     int      `yaml:"priority"`       // Handler priority (higher wins)
}

// TypeConversionConfig represents configuration for type-based field conversions.
type TypeConversionConfig struct {
	ProtoType           string `yaml:"proto_type"`            // Proto type to match (e.g., "string", "int32")
	DBType              string `yaml:"db_type"`               // DB type to match (e.g., "pgtype.UUID", "pgtype.Text")
	IsOptional          bool   `yaml:"is_optional"`           // Match optional fields
	ToProtoExpr         string `yaml:"to_proto_expr"`         // Conversion expression template (placeholders: {dbField}, {protoField}, {variable})
	ToDBExpr            string `yaml:"to_db_expr"`            // Conversion expression template (placeholders: {dbField}, {protoField}, {variable})
	Priority            int    `yaml:"priority"`              // Handler priority
	MatchFieldPattern   string `yaml:"match_field_pattern"`   // Regex pattern for field names
	MatchMessagePattern string `yaml:"match_message_pattern"` // Regex pattern for message names
	PointerStrategy     string `yaml:"pointer_strategy"`      // Pointer handling: strict, lenient, omit
	Alias               string `yaml:"alias"`                 // Reference to a type alias
}

// HelperFunctionConfig represents configuration for custom helper functions.
type HelperFunctionConfig struct {
	Name      string `yaml:"name"`      // Function name
	Signature string `yaml:"signature"` // Function signature
	Body      string `yaml:"body"`      // Function body
}

// ResponsePatternsConfig represents configuration for response field patterns.
type ResponsePatternsConfig struct {
	DataField      string   `yaml:"data_field"`      // Field name for data (default: "data")
	TotalField     string   `yaml:"total_field"`     // Field name for total (default: "total")
	PageField      string   `yaml:"page_field"`      // Field name for page (default: "page")
	LimitField     string   `yaml:"limit_field"`     // Field name for limit (default: "limit")
	ResponseSuffix string   `yaml:"response_suffix"` // Suffix for response messages (default: "Response")
	SkipFields     []string `yaml:"skip_fields"`     // Field names to skip in response helpers
}

// TypeAliasConfig represents configuration for reusable type conversions.
type TypeAliasConfig struct {
	ProtoType   string `yaml:"proto_type"`    // Proto type
	DBType      string `yaml:"db_type"`       // DB type
	IsOptional  bool   `yaml:"is_optional"`   // Optional flag
	ToProtoExpr string `yaml:"to_proto_expr"` // ToProto expression
	ToDBExpr    string `yaml:"to_db_expr"`    // ToDB expression
}

// PointerSettingsConfig represents configuration for pointer handling strategies.
type PointerSettingsConfig struct {
	DefaultStrategy string            `yaml:"default_strategy"` // Default strategy: strict, lenient, omit
	FieldStrategies map[string]string `yaml:"field_strategies"` // Per-field strategy overrides
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

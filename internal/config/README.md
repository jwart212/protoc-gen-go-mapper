# config

Package config handles loading and validating the mapper.yaml configuration file.

## Overview

The config package is responsible for reading the mapper.yaml configuration file, parsing it, and validating all fields before any downstream processing begins. This ensures fail-fast behavior for configuration errors.

## Configuration Format

```yaml
version: v1
database: sqlc
package:
  proto: internal/gen
  db: internal/postgres
```

### Fields

- **version**: Configuration version (must be "v1")
- **database**: Database backend (sqlc, pgx, or database_sql)
- **package.proto**: Go package path for generated protobuf code
- **package.db**: Go package path for database models

## Functions

### Load

Load reads and parses the mapper.yaml configuration from the given path:

```go
func Load(path string) (*Config, error)
```

### Validate

Validate checks that the configuration is valid:

```go
func Validate(cfg *Config) error
```

Validation rules:
- version must be "v1"
- database must be one of: sqlc, pgx, database_sql
- package.proto must be non-empty
- package.db must be non-empty

Any validation failure returns ErrInvalidConfig wrapped with the offending field name.

## Usage Example

```go
import "gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/internal/config"

cfg, err := config.Load("mapper.yaml")
if err != nil {
    // Handle error (will be ErrInvalidConfig for validation failures)
}
```

## Error Handling

All validation errors wrap ErrInvalidConfig with field-specific context:

```go
// Missing version
err := fmt.Errorf("validating config: version: %w", errors.ErrInvalidConfig)

// Invalid database value
err := fmt.Errorf("validating config: database %s: %w", dbValue, errors.ErrInvalidConfig)
```

## Design Decisions

- **Fail-fast validation**: Configuration is validated before any parsing begins to catch errors early.
- **Field-specific errors**: Each validation error includes the field name for clear error messages.
- **Strict versioning**: Only "v1" is supported to prevent silent mismatches.
- **Explicit database backends**: Only known database backends are accepted to catch typos.

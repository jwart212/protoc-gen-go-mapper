# protoc-gen-go-mapper

A protoc plugin that generates type-safe mapping functions between protobuf messages and database models.

## Overview

protoc-gen-go-mapper is a compiler-style plugin for protoc that generates bidirectional mapping functions between protobuf messages and database models (sqlc, pgx, database/sql). It uses a converter registry with priority-based resolution to handle type conversions deterministically.

## Features

- **Type-safe mappings**: Generated functions ensure compile-time type safety
- **Bidirectional**: Generates both Proto→DB and DB→Proto functions
- **Deterministic output**: Field order preserved by proto field numbers
- **Build-time validation**: All converter resolution happens during generation
- **Extensible**: Easy to add custom converters via the registry
- **No runtime reflection**: Pure Go code generation

## Installation

```bash
go install github.com/jwart212/protoc-gen-go-mapper/cmd/protoc-gen-go-mapper@latest
```

## Usage

1. Create a `mapper.yaml` configuration file:

```yaml
version: v1
database: sqlc
package:
  proto: internal/gen
  db: internal/postgres
```

2. Run protoc with the plugin:

```bash
protoc --go_out=. --go-mapper_out=. your.proto
```

3. Use the generated functions:

```go
// Convert DB model to protobuf
protoUser := ToProtoUser(dbUser)

// Convert protobuf to DB model
dbUser := ToDBUser(protoUser)
```

## Architecture

The plugin follows a compiler-style architecture:

1. **Parser** (`internal/parser/proto`): Converts proto descriptors to schema model
2. **Registry** (`internal/registry`): Manages converters with priority-based resolution
3. **Graph** (`internal/graph`): Builds mapping graphs from schema
4. **Generator** (`internal/generator`): Emits Go code using templates
5. **Plugin** (`internal/plugin`): Orchestrates the pipeline

## Supported Type Conversions

- Scalar types (string, int, bool, etc.)
- UUID ↔ string
- Timestamp ↔ time.Time
- Decimal ↔ string
- Enum ↔ string
- Nullable types (sql.Null*)
- Slice types ([]T ↔ []T)

## Configuration

The `mapper.yaml` file supports:

- `version`: Config version (must be "v1")
- `database`: Database type (sqlc, pgx, database_sql)
- `package.proto`: Go package for generated protobuf code
- `package.db`: Go package for database models

## Development

Run tests:

```bash
go test ./...
```

## License

MIT

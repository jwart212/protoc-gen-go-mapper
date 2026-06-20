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

### From GitHub

```bash
go install github.com/jwart212/protoc-gen-go-mapper/cmd/protoc-gen-go-mapper@latest
```

### From Source

```bash
git clone https://github.com/jwart212/protoc-gen-go-mapper.git
cd protoc-gen-go-mapper
go install ./cmd/protoc-gen-go-mapper
```

## Usage

### Basic Setup

1. Create a `mapper.yaml` configuration file:

```yaml
version: v1
database: sqlc
db_package: your-project/internal/db
package:
  proto: internal/gen
  db: internal/db
```

2. Define your protobuf message:

```protobuf
syntax = "proto3";

package user;

option go_package = "your-project/gen;gen";

message User {
  int32 id = 1;
  string name = 2;
  string email = 3;
  optional int32 age = 4;
}
```

3. Run protoc with the plugin:

```bash
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-mapper_out=. \
  --go-mapper_opt=paths=source_relative \
  user.proto
```

4. Use the generated functions:

```go
// Convert DB model to protobuf
protoUser := ToProtoUser(dbUser)

// Convert protobuf to DB model
dbUser := ToDBUser(protoUser)
```

### Advanced Configuration

#### Database Types

The plugin supports three database backends:

**SQLC (PostgreSQL with pgtype):**
```yaml
database: sqlc
db_package: your-project/internal/postgres/sqlc
```

**PGX:**
```yaml
database: pgx
db_package: your-project/internal/pgx
```

**database/sql:**
```yaml
database: database_sql
db_package: your-project/internal/db
```

#### Configuration Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `version` | string | Yes | Config version (must be "v1") |
| `database` | string | Yes | Database type: `sqlc`, `pgx`, or `database_sql` |
| `db_package` | string | Yes | Go package path for database models |
| `package.proto` | string | Yes | Go package for generated protobuf code |
| `package.db` | string | Yes | Go package for database models |

## Architecture

The plugin follows a compiler-style architecture:

1. **Parser** (`internal/parser/proto`): Converts proto descriptors to schema model
2. **Registry** (`internal/registry`): Manages converters with priority-based resolution
3. **Graph** (`internal/graph`): Builds mapping graphs from schema
4. **Generator** (`internal/generator`): Emits Go code using templates
5. **Plugin** (`internal/plugin`): Orchestrates the pipeline

## Supported Type Conversions

### Scalar Types
- Basic Go types: `string`, `int32`, `int64`, `bool`, `float32`, `float64`
- Direct mapping between proto and DB scalar types

### UUID
- `uuid.UUID` ↔ `string`
- `pgtype.UUID` ↔ `string` (SQLC/PGX)
- Handles nullable UUID fields

### Timestamp
- `time.Time` ↔ `google.protobuf.Timestamp`
- `pgtype.Timestamptz` ↔ `google.protobuf.Timestamp` (SQLC)
- `pgtype.Timestamp` ↔ `google.protobuf.Timestamp` (PGX)
- `sql.NullTime` ↔ `*google.protobuf.Timestamp`

### Decimal
- `decimal.Decimal` ↔ `string`
- `pgtype.Numeric` ↔ `string` (SQLC/PGX)

### Enum
- Proto enum ↔ SQLC enum (string)
- Proto enum ↔ string
- Handles zero-value mapping for unknown DB values

### Nullable Types
- `sql.NullString` ↔ `optional string`
- `sql.NullInt64` ↔ `optional int64`
- `sql.NullBool` ↔ `optional bool`
- `sql.NullFloat64` ↔ `optional double`
- `sql.NullTime` ↔ `optional google.protobuf.Timestamp`

### Slice Types
- `[]T` ↔ `[]T` with element-level conversion
- Supports nested nullable types in slices

## Examples

### Simple User Mapping

**Protobuf:**
```protobuf
message User {
  int32 id = 1;
  string name = 2;
  string email = 3;
  optional int32 age = 4;
}
```

**Generated Code:**
```go
func ToProtoUser(src db.User) *User {
  return &User{
    Id:    src.ID,
    Name:  src.Name,
    Email: src.Email,
    Age:   newInt32(src.Age.Int32, src.Age.Valid),
  }
}

func ToDBUser(src *User) db.User {
  return db.User{
    ID:   src.Id,
    Name: src.Name,
    Email: src.Email,
    Age:  newNullInt32(src.Age),
  }
}
```

### UUID and Timestamp Mapping

**Protobuf:**
```protobuf
message Product {
  string id = 1;
  string name = 2;
  google.protobuf.Timestamp created_at = 3;
}
```

**Generated Code (SQLC):**
```go
func ToProtoProduct(src db.Product) *Product {
  return &Product{
    Id:        newStringFromUUIDNonPtr(src.ID),
    Name:      src.Name,
    CreatedAt: newTimestampFromTimestamptz(src.CreatedAt),
  }
}

func ToDBProduct(src *Product) db.Product {
  return db.Product{
    ID:        newUUIDFromString(src.Id),
    Name:      src.Name,
    CreatedAt: newTimestamptzFromTimestamp(src.CreatedAt),
  }
}
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Building

```bash
# Build the plugin
go build -o bin/protoc-gen-go-mapper ./cmd/protoc-gen-go-mapper

# Install to GOPATH
go install ./cmd/protoc-gen-go-mapper
```

### Project Structure

```
protoc-gen-go-mapper/
├── cmd/protoc-gen-go-mapper/    # Main plugin entry point
├── internal/
│   ├── config/                   # Configuration parsing
│   ├── generator/                # Code generation
│   ├── graph/                    # Mapping graph construction
│   ├── parser/proto/             # Proto descriptor parsing
│   ├── registry/                 # Converter registry
│   ├── resolver/                 # Type resolution
│   └── schema/                   # Schema model
├── pkg/
│   ├── converter/                # Converter interface
│   ├── errors/                   # Error types
│   ├── naming/                   # Naming conventions
│   └── types/                    # Type system
└── examples/                     # Usage examples
```

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Run `gofmt` on all code
- Add tests for new features
- Update documentation as needed

## Support

If you find this project useful, consider supporting its development:

[☕ Buy me a coffee](https://www.buymeacoffee.com/ekowdd89)

## License

MIT License - see LICENSE file for details.

## GitHub Repository

https://github.com/jwart212/protoc-gen-go-mapper

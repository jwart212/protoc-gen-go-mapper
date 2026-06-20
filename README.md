# protoc-gen-go-mapper

[![Go Version](https://img.shields.io/badge/Go-1.25.0+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

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
- **Zero-config mode**: Automatic type detection with generic converters
- **Self-contained**: Generated code includes inline converters, no external dependencies

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Dependencies](#dependencies)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Advanced Configuration](#advanced-configuration)
- [Examples](#examples)
- [Development](#development)
- [Contributing](#contributing)
- [Support](#support)
- [License](#license)

## Dependencies

### Required Dependencies

- **gopkg.in/yaml.v3** v3.0.1 - YAML configuration parsing
- **google.golang.org/protobuf** v1.36.11 - Protocol Buffers support
- **github.com/jackc/pgx/v5** v5.10.0 - PostgreSQL driver and types (for SQLC/PGX support)

### Database-Specific Dependencies

The generated code works with the following database libraries (not required by the plugin itself):

- **SQLC** - SQL query code generator (for SQLC mode)
- **github.com/jackc/pgx/v5** - PostgreSQL driver (for PGX mode)
- **database/sql** - Standard database SQL interface (for database_sql mode)

### Indirect Dependencies

- **github.com/google/uuid** v1.6.0 - UUID parsing and generation

### Go Version

- **Go 1.25.0** or higher

### Installation

To install the required dependencies:

```bash
go mod download
```

To update dependencies:

```bash
go get -u ./...
go mod tidy
```

## Installation

### Quick Install

```bash
go install github.com/jwart212/protoc-gen-go-mapper/cmd/protoc-gen-go-mapper@latest
```

### From Source

```bash
git clone https://github.com/jwart212/protoc-gen-go-mapper.git
cd protoc-gen-go-mapper
go install ./cmd/protoc-gen-go-mapper
```

## Quick Start

### 1. Create a simple mapper.yaml

```yaml
version: v1
database: sqlc
db_package: your-project/internal/postgres/sqlc
package:
  proto: internal/gen
  db: internal/postgres
type_mappings:
  User: DbUser
messages:
  - User
```

### 2. Define your protobuf message

```protobuf
syntax = "proto3";

package user;

option go_package = "your-project/gen;gen";

message User {
  string id = 1;
  string name = 2;
  string email = 3;
}
```

### 3. Run protoc

```bash
protoc \
  --proto_path=internal/proto \
  --go_out=. \
  --go-grpc_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  --plugin=protoc-gen-mapper=protoc-gen-go-mapper \
  --mapper_out=. \
  --mapper_opt=paths=source_relative,mapper_config=internal/proto/mapper.yaml \
  internal/proto/user.proto
```

### 4. Use the generated functions

```go
// Convert DB model to protobuf
protoUser := ToProtoUser(dbUser)

// Convert protobuf to DB model
dbUser := ToDBUser(protoUser)
```

## Usage

### Advanced Configuration

#### mapper.yaml Configuration Reference

The `mapper.yaml` file controls how the plugin generates mapping functions. Here's a comprehensive reference of all available options:

```yaml
# Mapper configuration for protoc-gen-go-mapper
version: v1
database: sqlc
db_package: github.com/your-project/internal/postgres/sqlc
package:
  proto: internal/gen
  db: internal/postgres

# Type mappings between proto messages and database models
type_mappings:
  User: DbUser
  CreateUserRequest: CreateUserParams
  UpdateUserRequest: UpdateUserParams

# Response type mappings for list responses (proto message -> SQLC Row type)
response_type_mappings:
  ListUserResponse: ListUsersRow
  ListProductResponse: ListProductsRow

# Type aliases for reuse (optional, for custom conversions)
type_aliases:
  UUIDField:
    proto_type: string
    db_type: pgtype.UUID
    is_optional: false
    to_proto_expr: "newStringFromUUIDNonPtr({variable}.{dbField})"
    to_db_expr: "newUUIDFromString({protoField})"

# Type-based field conversions (optional - remove to use generic converters)
# When omitted, the plugin uses built-in generic converters for common types
type_conversions:
  # Non-nullable UUID fields with pattern matching
  - proto_type: string
    db_type: pgtype.UUID
    is_optional: false
    match_field_pattern: "^.*Id$"
    to_proto_expr: "newStringFromUUIDNonPtr({variable}.{dbField})"
    to_db_expr: "newUUIDFromString({protoField})"
    priority: 90

# Response field patterns (configurable)
response_patterns:
  data_field: "data"
  total_field: "total"
  page_field: "page"
  limit_field: "limit"
  response_suffix: "Response"
  skip_fields: ["DeletedAt", "DeletedBy"]

# Pointer handling strategies (optional)
pointer_settings:
  default_strategy: "strict"
  field_strategies:
    CreatedAt: "lenient"
    UpdatedAt: "lenient"
    DeletedAt: "omit"

# Messages to generate mappers for
messages:
  - User
  - CreateUserRequest
  - UpdateUserRequest

# Field handlers for special cases
field_handlers:
  # Type assertion for interface{} fields
  - name: path_type_assertion
    type: type_assertion
    match_field: path
    match_db_types:
      - ListCategoriesTreeRow
    assert_type: "[]string"
    priority: 50

  # Default value for fields
  - name: children_default_value
    type: default_value
    match_field: children
    match_db_types:
      - ListCategoriesTreeRow
    default_value: "[]*CategoryTreeNode{}"
    priority: 75
```

#### Configuration Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `version` | string | ✅ Yes | Config version (must be "v1") |
| `database` | string | ✅ Yes | Database type: `sqlc`, `pgx`, or `database_sql` |
| `db_package` | string | ✅ Yes | Go package path for database models |
| `package.proto` | string | ✅ Yes | Go package for generated protobuf code |
| `package.db` | string | ✅ Yes | Go package for database models |
| `type_mappings` | object | ❌ No | Custom type mappings between proto messages and DB models |
| `response_type_mappings` | object | ❌ No | Mappings for response messages to SQLC Row types |
| `type_aliases` | object | ❌ No | Reusable type conversion definitions |
| `type_conversions` | array | ❌ No | Custom type conversion rules (optional - generic converters used when omitted) |
| `response_patterns` | object | ❌ No | Configuration for response helper generation |
| `pointer_settings` | object | ❌ No | Pointer handling strategies |
| `messages` | array | ❌ No | List of proto messages to generate mappers for |
| `field_handlers` | array | ❌ No | Special field handling rules |

#### Generic Converters (Zero-Config Mode)

When `type_conversions` is omitted or empty, the plugin automatically uses built-in generic converters for common types:

**Supported Generic Conversions:**
- `string` ↔ `pgtype.UUID` (UUID fields)
- `string` ↔ `pgtype.Text` (Text fields)
- `google.protobuf.Timestamp` ↔ `pgtype.Timestamptz` (Timestamp fields)
- Optional fields handled automatically (e.g., `*string` for nullable UUID)

**Example Zero-Config mapper.yaml:**
```yaml
version: v1
database: sqlc
db_package: github.com/your-project/internal/postgres/sqlc
package:
  proto: internal/gen
  db: internal/postgres
type_mappings:
  User: DbUser
messages:
  - User
```

This minimal configuration will generate correct mappings for UUID, Timestamp, and Text fields automatically.

#### Custom Type Conversions

For custom type conversions not covered by generic converters, use `type_conversions`:

```yaml
type_conversions:
  # Custom decimal conversion
  - proto_type: string
    db_type: pgtype.Numeric
    is_optional: false
    to_proto_expr: "decimalToString({variable}.{dbField})"
    to_db_expr: "stringToDecimal({protoField})"
    priority: 90

  # Pattern-based matching (regex)
  - proto_type: string
    db_type: pgtype.UUID
    is_optional: false
    match_field_pattern: "^.*Id$"  # Matches fields ending with "Id"
    to_proto_expr: "newStringFromUUIDNonPtr({variable}.{dbField})"
    to_db_expr: "newUUIDFromString({protoField})"
    priority: 90
```

**Type Conversion Options:**
- `proto_type`: Protobuf field type
- `db_type`: Database field type
- `is_optional`: Whether the field is optional/nullable
- `match_field_pattern`: Regex pattern to match field names (case-insensitive)
- `to_proto_expr`: Expression for DB→Proto conversion
- `to_db_expr`: Expression for Proto→DB conversion
- `priority`: Priority for matching (higher = more specific)

#### Field Handlers

Field handlers handle special cases like type assertions and default values:

```yaml
field_handlers:
  # Type assertion for interface{} fields
  - name: path_type_assertion
    type: type_assertion
    match_field: path
    match_db_types:
      - ListCategoriesTreeRow
    assert_type: "[]string"
    priority: 50

  # Default value for fields
  - name: children_default_value
    type: default_value
    match_field: children
    match_db_types:
      - ListCategoriesTreeRow
    default_value: "[]*CategoryTreeNode{}"
    priority: 75

  # Skip fields
  - name: skip_internal_fields
    type: skip
    match_field_pattern: "^internal_"
    priority: 100
```

**Field Handler Types:**
- `type_assertion`: Asserts a type for interface{} fields
- `default_value`: Sets a default value for a field
- `skip`: Skips the field during mapping

#### Response Patterns

Configure response helper generation for list responses:

```yaml
response_patterns:
  data_field: "data"           # Field name for data array
  total_field: "total"         # Field name for total count
  page_field: "page"           # Field name for page number
  limit_field: "limit"         # Field name for page limit
  response_suffix: "Response"  # Suffix for response messages
  skip_fields:                 # Fields to skip in response
    - DeletedAt
    - DeletedBy
```

#### Pointer Settings

Control how nullable fields are handled:

```yaml
pointer_settings:
  default_strategy: "strict"   # Default strategy: strict, lenient, or omit
  field_strategies:
    CreatedAt: "lenient"       # Allow null for optional fields
    UpdatedAt: "lenient"
    DeletedAt: "omit"          # Skip field if null
```

**Pointer Strategies:**
- `strict`: Require non-null values (default)
- `lenient`: Allow null values for optional fields
- `omit`: Skip field if value is null

#### Example mapper.yaml (Advanced Example)

Here's a detailed breakdown of the example configuration from `examples/advanced/internal/proto/mapper.yaml`:

```yaml
# Mapper configuration for protoc-gen-go-mapper
version: v1
database: sqlc
db_package: github.com/jwart212/protoc-gen-go-mapper/examples/advanced/internal/postgres/sqlc
package:
  proto: internal/gen
  db: internal/postgres

# Type mappings between proto messages and database models
type_mappings:
  UOM: SchmPosUom                                    # Unit of Measure
  ItemCategory: SchmPosItemCategory                  # Item Category
  ListsItemCategory: ListsItemCategoriesRow         # List response
  CategoryTreeNode: ListsItemCategoriesTreeRow       # Tree structure
  CreateUOMRequest: CreateUOMParams                 # Create parameters
  UpdateUOMRequest: UpdatedUOMParams                 # Update parameters
  CreateItemCategoryRequest: CreateItemCategoriesParams
  UpdateItemCategoryRequest: UpdateItemCategoriesParams
  DeleteItemCategoryRequest: DeleteItemCategoriesParams

# Response type mappings for list responses
response_type_mappings:
  ListItemCategoryResponse: ListsItemCategoriesRow
  ListItemCategoryTreeResponse: ListsItemCategoriesTreeRow

# Type aliases for reuse (optional, for custom conversions)
type_aliases:
  UUIDField:
    proto_type: string
    db_type: pgtype.UUID
    is_optional: false
    to_proto_expr: "newStringFromUUIDNonPtr({variable}.{dbField})"
    to_db_expr: "newUUIDFromString({protoField})"
  
  OptionalUUID:
    proto_type: string
    db_type: pgtype.UUID
    is_optional: true
    to_proto_expr: "newStringFromUUID({variable}.{dbField})"
    to_db_expr: "newUUIDFromString({protoField})"

# Type-based field conversions (optional - commented out to use generic converters)
# type_conversions:
#   # Non-nullable UUID fields with pattern matching
#   - proto_type: string
#     db_type: pgtype.UUID
#     is_optional: false
#     match_field_pattern: "^.*Id$"
#     to_proto_expr: "newStringFromUUIDNonPtr({variable}.{dbField})"
#     to_db_expr: "newUUIDFromString({protoField})"
#     priority: 90
#   
#   # Optional UUID fields with pattern matching
#   - proto_type: string
#     db_type: pgtype.UUID
#     is_optional: true
#     match_field_pattern: "parentid|parent_id"
#     to_proto_expr: "newStringFromUUID({variable}.{dbField})"
#     to_db_expr: "newUUIDFromString({protoField})"
#     priority: 90
#   
#   # Text fields (Description)
#   - proto_type: string
#     db_type: pgtype.Text
#     is_optional: true
#     to_proto_expr: "newString({variable}.{dbField}.String, {variable}.{dbField}.Valid)"
#     to_db_expr: "newText({protoField})"
#     priority: 90
#   
#   # Non-optional timestamp fields (CreatedAt, UpdatedAt)
#   - proto_type: google.protobuf.Timestamp
#     db_type: pgtype.Timestamptz
#     is_optional: false
#     to_proto_expr: "newTimestampFromTimestamptz({variable}.{dbField})"
#     to_db_expr: "newTimestamptzFromTimestamp({protoField})"
#     priority: 90
#   
#   # Optional timestamp fields (DeletedAt)
#   - proto_type: google.protobuf.Timestamp
#     db_type: pgtype.Timestamptz
#     is_optional: true
#     to_proto_expr: "newTimestampFromTimestamptz({variable}.{dbField})"
#     to_db_expr: "newTimestamptzFromTimestamp({protoField})"
#     priority: 90

# Response field patterns (configurable)
response_patterns:
  data_field: "data"
  total_field: "total"
  page_field: "page"
  limit_field: "limit"
  response_suffix: "Response"
  skip_fields: ["DeletedAt", "DeletedBy"]  # Skip soft delete fields in responses

# Pointer handling strategies (optional)
pointer_settings:
  default_strategy: "strict"
  field_strategies:
    CreatedAt: "lenient"   # Allow null for optional timestamp fields
    UpdatedAt: "lenient"
    DeletedAt: "omit"      # Skip DeletedAt if null

# Messages to generate mappers for
messages:
  - UOM
  - ItemCategory
  - ListsItemCategory
  - CategoryTreeNode
  - CreateUOMRequest
  - UpdateUOMRequest
  - CreateItemCategoryRequest
  - UpdateItemCategoryRequest
  - DeleteItemCategoryRequest

# Field handlers for special cases
field_handlers:
  # Type assertion for Path field (interface{} to []string)
  # This handles the special case where SQLC returns interface{} for array fields
  - name: path_type_assertion
    type: type_assertion
    match_field: path
    match_db_types:
      - ListsItemCategoriesTreeRow
    assert_type: "[]string"
    priority: 50
  
  # Default value for Children field in CategoryTreeNode
  # Since the DB row doesn't have a Children field, we set it to empty slice
  - name: children_default_value
    type: default_value
    match_field: children
    match_db_types:
      - ListsItemCategoriesTreeRow
    default_value: "[]*CategoryTreeNode{}"
    priority: 75
```

**Key Points from This Example:**

1. **Generic Converters Enabled**: The `type_conversions` section is commented out, so the plugin uses built-in generic converters for UUID, Timestamp, and Text fields automatically.

2. **Type Aliases**: Defined but not used in this example (would be used if `type_conversions` were active).

3. **Response Patterns**: Configured to skip soft delete fields (`DeletedAt`, `DeletedBy`) in list responses.

4. **Field Handlers**: Two special cases:
   - `path_type_assertion`: Converts `interface{}` to `[]string` for the Path field in tree structures
   - `children_default_value`: Sets empty slice for Children field since it doesn't exist in the DB row

5. **Pointer Settings**: Configured to be lenient with timestamp fields and omit soft delete fields when null.

**Generated Code with Generic Converters:**

When using this configuration with generic converters enabled, the generated code includes inline converter functions:

```go
// Generic converter functions (inlined in generated file)
func ConvertUUID[T string | *string](v pgtype.UUID) T {
	if v.Valid {
		s := uuid.UUID(v.Bytes).String()
		var t T
		switch any(t).(type) {
		case string:
			return any(s).(T)
		case *string:
			return any(&s).(T)
		}
	}
	var zero T
	return zero
}

func ConvertTimestamp[T *timestamppb.Timestamp](v pgtype.Timestamptz) T {
	if v.Valid {
		return timestamppb.New(v.Time)
	}
	var zero T
	return zero
}

func ConvertText[T string | *string](v pgtype.Text) T {
	if v.Valid {
		var t T
		switch any(t).(type) {
		case string:
			return any(v.String).(T)
		case *string:
			return any(&v.String).(T)
		}
	}
	var zero T
	return zero
}

// Usage in generated functions
func ToProtoItemCategory(src sqlc.SchmPosItemCategory) *ItemCategory {
	return &ItemCategory{
		Id: ConvertUUID[string](src.ID),
		TenantId: ConvertUUID[string](src.TenantID),
		ParentId: ConvertUUID[*string](src.ParentID),
		Description: ConvertText[*string](src.Description),
		CreatedAt: ConvertTimestamp[*timestamppb.Timestamp](src.CreatedAt),
		UpdatedAt: ConvertTimestamp[*timestamppb.Timestamp](src.UpdatedAt),
	}
}
```

This approach makes the generated code self-contained without requiring external dependencies on the plugin's internal packages.

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
| `type_mappings` | object | No | Custom type mappings between proto messages and DB models |
| `messages` | array | No | List of proto messages to generate mappers for |

#### Type Mappings

The `type_mappings` section allows you to specify custom mappings between protobuf message types and database model types:

```yaml
type_mappings:
  User: DbUser                    # Proto message -> DB model
  CreateUserRequest: CreateUserParams  # Proto request -> DB params
  UpdateUserRequest: UpdateUserParams
```

#### Messages

The `messages` section specifies which protobuf messages should have mapping functions generated:

```yaml
messages:
  - User
  - CreateUserRequest
  - UpdateUserRequest
```

If `messages` is not specified, all messages in the proto file will be processed.

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

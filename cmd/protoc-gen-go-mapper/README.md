# protoc-gen-go-mapper

protoc-gen-go-mapper is a protoc plugin that generates type-safe mapping functions between protobuf messages and database models.

## Installation

```bash
go install gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/cmd/protoc-gen-go-mapper@latest
```

## Usage

Add the plugin to your protoc invocation:

```bash
protoc --go_out=. --go-mapper_out=. your.proto
```

## Configuration

Create a `mapper.yaml` file in your project root:

```yaml
version: v1
database: sqlc
package:
  proto: internal/gen
  db: internal/postgres
```

## Generated Code

For each protobuf message, the plugin generates:

```go
func ToProtoUser(src db.User) *pb.User
func ToDBUser(src *pb.User) db.User
```

## Supported Databases

- sqlc
- pgx
- database/sql

## Supported Type Conversions

- Scalar types (string, int, bool, etc.)
- UUID ↔ string
- Timestamp ↔ time.Time
- Decimal ↔ string
- Enum ↔ string
- Nullable types
- Slice types

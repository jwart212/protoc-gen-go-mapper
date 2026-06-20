# Simple Example

A simple example demonstrating protoc-gen-go-mapper with a single table and basic types.

## Schema

PostgreSQL table with basic types:
- `id`: SERIAL (auto-increment integer)
- `name`: TEXT (string)
- `email`: TEXT (string)
- `age`: INTEGER (optional)
- `active`: BOOLEAN (default true)

## Setup

1. Generate sqlc code:
```bash
cd examples/simple
sqlc generate
```

2. Generate protobuf code:
```bash
protoc --go_out=. --go_opt=paths=source_relative user.proto
```

3. Generate mapper code:
```bash
protoc --go-mapper_out=. --go-mapper_opt=mapper.example.yaml user.proto
```

## Generated Functions

For the `User` message, the plugin generates:
- `ToProtoUser(db.User) *gen.User` - Convert DB model to protobuf
- `ToDBUser(*gen.User) db.User` - Convert protobuf to DB model

## Usage

```go
// Convert DB model to protobuf
protoUser := ToProtoUser(dbUser)

// Convert protobuf to DB model
dbUser := ToDBUser(protoUser)
```

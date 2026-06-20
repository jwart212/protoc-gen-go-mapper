# Advanced Example: Real-World gRPC Service with protoc-gen-go-mapper

This example demonstrates a production-grade gRPC service using `protoc-gen-go-mapper` for automatic type conversions between Protobuf and SQLC models.

## Overview

This example simulates a Point of Sale (POS) system with:
- **UOM (Unit of Measure)** management
- **Item Category** management
- **SQLC** for type-safe database queries
- **protoc-gen-go-mapper** for automatic proto-to-DB conversions

## Key Features Demonstrated

### 1. UUID Handling
- Automatic conversion between `string` (protobuf) and `pgtype.UUID` (SQLC/PostgreSQL)
- Supports nullable UUID fields with `uuid.NullUUID`

### 2. Timestamp Handling
- Automatic conversion between `google.protobuf.Timestamp` and `pgtype.Timestamptz`
- Supports nullable timestamps for soft delete patterns

### 3. Nullable Fields
- Proto3 `optional` fields correctly mapped to nullable SQL types
- `sql.NullString`, `sql.NullTime`, etc. handled automatically

### 4. Field Name Mapping
- Automatic snake_case ↔ camelCase conversion
- DB: `tenant_id`, `created_at`, `deleted_at`
- Proto: `tenantId`, `createdAt`, `deletedAt`

## Project Structure

```
examples/advanced/
├── cmd/
│   └── main.go              # Demo application
├── internal/
│   ├── gen/
│   │   ├── uompb/           # Generated UOM mappers
│   │   └── item_categoriespb/ # Generated ItemCategory mappers
│   ├── postgres/
│   │   ├── schema.sql       # PostgreSQL schema
│   │   ├── query/
│   │   │   └── query.sql    # SQLC queries
│   │   └── sqlc/            # SQLC generated code
│   └── proto/
│       ├── uom.proto        # gRPC service definitions
│       ├── item_category.proto
│       ├── uom_model.proto  # Model-only proto for mapping
│       ├── item_category_model.proto
│       └── mapper.yaml      # Mapper configuration
```

## Generated Mapper Functions

For each model, `protoc-gen-go-mapper` generates:

```go
// Convert protobuf to SQLC model
func ToDBUOM(src *UOM) sqlc.SchmPosUom

// Convert SQLC model to protobuf
func ToProtoUOM(src sqlc.SchmPosUom) *UOM
```

## Usage Example

```go
// Create a protobuf UOM
protoUOM := uompb.UOM{
    Id:          uuid.New().String(),
    TenantId:    uuid.New().String(),
    Code:        "KG",
    Name:        "Kilogram",
    Symbol:      "kg",
    Description: "Metric unit of mass",
    CreatedAt:   timestamppb.New(time.Now()),
    UpdatedAt:   timestamppb.New(time.Now()),
}

// Convert to DB model automatically
dbUOM := uompb.ToDBUOM(&protoUOM)

// Use with SQLC
result, err := queries.CreateUOM(ctx, dbUOM)

// Convert back to protobuf
protoUOM2 := uompb.ToProtoUOM(result)
```

## Database Schema

The schema uses PostgreSQL with:
- UUID primary keys
- TIMESTAMPTZ for timestamps
- Soft delete pattern (deleted_at, deleted_by)
- Tenant isolation (tenant_id)

## gRPC Services

### UOMService
- Create, Update, Delete, Restore
- Get, List, Exists
- BatchDelete, BatchRestore

### ItemCategoryService
- Create, Update, Delete, Restore
- Get, List, Tree (hierarchical)
- Parent-child relationships

## Running the Example

### Prerequisites
- Go 1.24+
- PostgreSQL
- protoc compiler
- protoc-gen-go
- protoc-gen-go-grpc
- protoc-gen-go-mapper (this plugin)
- sqlc

### Generate Code

```bash
# Generate protobuf code
protoc --go_out=internal/gen/internal/proto \
       --go_opt=paths=source_relative \
       --go-grpc_out=internal/gen/internal/proto \
       --go-grpc_opt=paths=source_relative \
       internal/proto/*.proto

# Generate mapper code
protoc --go-mapper_out=internal/gen/uompb \
       --go-mapper_opt=mapper.yaml \
       --plugin=protoc-gen-go-mapper=../../protoc-gen-go-mapper.exe \
       internal/proto/uom_model.proto

protoc --go-mapper_out=internal/gen/item_categoriespb \
       --go-mapper_opt=mapper.yaml \
       --plugin=protoc-gen-go-mapper=../../protoc-gen-go-mapper.exe \
       internal/proto/item_category_model.proto

# Generate SQLC code
cd internal/postgres
sqlc generate
```

### Run Demo

```bash
cd cmd
go run main.go
```

## Benefits of protoc-gen-go-mapper

1. **No Boilerplate**: Eliminates manual conversion code
2. **Type Safety**: Compile-time type checking
3. **Maintainability**: Single source of truth (proto definitions)
4. **Consistency**: All conversions follow the same pattern
5. **Performance**: Direct field assignments, no reflection

## Real-World Use Case

This example demonstrates a typical microservice architecture where:
- gRPC API receives protobuf messages
- Service layer uses protobuf types
- Repository layer uses SQLC types
- Database uses PostgreSQL types

`protoc-gen-go-mapper` bridges the gap between these layers automatically.

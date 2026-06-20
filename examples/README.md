# Examples

This directory contains example projects demonstrating protoc-gen-go-mapper with sqlc, protobuf, and PostgreSQL.

## Examples Overview

### Simple Example (`simple/`)
A basic example with a single table and basic types.
- **Schema**: Single `users` table with id, name, email, age, active
- **Types**: SERIAL, TEXT, INTEGER, BOOLEAN
- **Features**: Basic CRUD operations
- **Complexity**: Beginner

### Medium Example (`medium/`)
A medium example with multiple tables and relationships.
- **Schema**: `authors`, `books`, `reviews` tables
- **Relationships**: One-to-many (authors → books, books → reviews)
- **Types**: TEXT, INTEGER, DATE, TIMESTAMP
- **Features**: Foreign keys, JOIN queries
- **Complexity**: Intermediate

### Complex Example (`complex/`)
An advanced example with nested messages, enums, and nullable types.
- **Schema**: `customers`, `products`, `orders`, `order_items` tables
- **Relationships**: Many-to-many through order_items
- **Types**: UUID, custom enum, DECIMAL, TIMESTAMP, nullable fields
- **Features**: UUID primary keys, custom enums, nested messages, repeated fields
- **Complexity**: Advanced

## Common Setup Steps

For each example, follow these steps:

1. **Generate sqlc code:**
```bash
cd examples/[simple|medium|complex]
sqlc generate
```

2. **Generate protobuf code:**
```bash
protoc --go_out=. --go_opt=paths=source_relative [user|library|ecommerce].proto
```

3. **Generate mapper code:**
```bash
protoc --go-mapper_out=. --go-mapper_opt=mapper.example.yaml [user|library|ecommerce].proto
```

## Type Conversions Demonstrated

| Example | Conversions |
|---------|-------------|
| Simple | Scalar types (int32, string, bool) |
| Medium | Scalar types + DATE |
| Complex | UUID, Enum, Timestamp, Decimal, Nullable, Repeated |

## Requirements

- PostgreSQL database
- sqlc CLI tool
- protoc compiler
- protoc-gen-go plugin
- protoc-gen-go-mapper plugin

## File Structure

Each example directory contains:
- `schema.sql`: PostgreSQL schema and queries
- `[name].proto`: Protobuf message definitions
- `sqlc.yaml`: sqlc configuration
- `mapper.example.yaml`: Mapper configuration (rename to mapper.yaml to use)
- `README.md`: Example-specific documentation


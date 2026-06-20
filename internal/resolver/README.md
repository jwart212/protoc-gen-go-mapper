# resolver

Package resolver maps protobuf types to database-specific types.

## Overview

The resolver package provides database-specific type mapping for sqlc, pgx, and database/sql.

## Types

### Resolver

Resolver maps protobuf types to database-specific types:

```go
type Resolver struct {
    database string
}

func New(database string) *Resolver
func (r *Resolver) Resolve(protoType types.TypeInfo) types.TypeInfo
```

## Supported Databases

- **sqlc**: Maps UUID → string, Timestamp → time.Time, Decimal → string
- **pgx**: Maps UUID → pgtype.UUID, Timestamp → pgtype.Timestamp, Decimal → pgtype.Numeric
- **database_sql**: Maps UUID → string, Timestamp → time.Time, Decimal → string

## Usage Example

```go
import "github.com/jwart212/protoc-gen-go-mapper/internal/resolver"

r := resolver.New("pgx")
dbType := r.Resolve(protoType)
```

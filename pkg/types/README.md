# types

Package types provides the internal type system used throughout the mapper.

## Overview

The types package defines the core type representation used for converter matching and code generation. It deliberately avoids raw strings outside the registry package to ensure type safety.

## Types

### TypeInfo

TypeInfo describes a single Go type as understood by the mapper:

```go
type TypeInfo struct {
    Package string
    Name    string
    IsPointer  bool
    IsSlice    bool
    IsEnum     bool
    IsNullable bool
    Kind       Kind
}
```

TypeInfo is a value type (passed by value) since it is small and immutable by convention.

### Kind

Kind classifies a TypeInfo for converter matching:

```go
const (
    KindScalar Kind = iota
    KindUUID
    KindTimestamp
    KindDecimal
    KindEnum
    KindNullable
    KindMessage
)
```

Kind implements String() for human-readable output in error messages and logging.

## Usage Example

```go
import "gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/types"

// Create a TypeInfo for a UUID field
uuidType := types.TypeInfo{
    Package: "github.com/google/uuid",
    Name:    "UUID",
    Kind:    types.KindUUID,
}

// Check the kind
if uuidType.Kind == types.KindUUID {
    // Handle UUID conversion
}
```

## Design Decisions

- **Value semantics**: TypeInfo is passed by value, not by pointer, following Effective Go's guidance for small immutable structs.
- **No raw strings**: Type identification uses Kind enum rather than string comparisons to avoid typos and enable compile-time checking.
- **String() method**: Kind implements String() to make error messages and logging human-readable without sacrificing type safety.

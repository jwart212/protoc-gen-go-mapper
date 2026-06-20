# schema

Package schema defines the internal schema model used to represent protobuf messages and enums.

## Overview

The schema package provides the data structures that represent parsed protobuf descriptors. These structures are used throughout the mapper pipeline after parsing and before type resolution.

## Types

### Model

Model represents the complete schema parsed from protobuf descriptors:

```go
type Model struct {
    Messages []*Message
    Enums    []*Enum
}
```

### Message

Message represents a protobuf message definition:

```go
type Message struct {
    Name   string
    Fields []*Field
}
```

### Field

Field represents a single field in a protobuf message:

```go
type Field struct {
    Name string
    ProtoType types.TypeInfo
    DBType    types.TypeInfo
    FieldNumber int32
    Repeated bool
    Optional bool
}
```

**Important**: FieldNumber preserves the original proto field declaration order. This is the canonical sort key for deterministic output. Never re-derive order from map iteration.

### Enum

Enum represents a protobuf enum definition:

```go
type Enum struct {
    Name   string
    Values []string
}
```

## Determinism

Field ordering in generated code MUST follow Field.FieldNumber ascending, which mirrors the original .proto declaration order. This is enforced through the schema model and validated by determinism tests.

## Usage Example

```go
import "github.com/jwart212/protoc-gen-go-mapper/internal/schema"

msg := &schema.Message{
    Name: "User",
    Fields: []*schema.Field{
        {
            Name:        "id",
            FieldNumber: 1,
        },
        {
            Name:        "name",
            FieldNumber: 2,
        },
    },
}
```

## Design Decisions

- **FieldNumber preservation**: FieldNumber is explicitly stored to ensure deterministic output order matching proto declaration.
- **Value types**: Field and Enum are value types where appropriate to avoid unnecessary pointer indirection.
- **Pointer slices**: Message.Fields is a slice of pointers to allow mutation during graph construction.

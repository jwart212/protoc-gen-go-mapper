# graph

Package graph builds the mapping graph that connects source and target fields with their converters.

## Overview

The graph package constructs the mapping graph used by the generator. It represents the complete mapping between two message types with field-to-field connections and their assigned converters.

## Types

### Mapper

Mapper represents the complete mapping between two message types:

```go
type Mapper struct {
    Source string
    Target string
    Fields []FieldMapping
}
```

### FieldMapping

FieldMapping connects individual fields with their assigned converter:

```go
type FieldMapping struct {
    SourceField string
    TargetField string
    Converter   Converter
}
```

## Functions

### NewMapper

NewMapper creates a new Mapper for the given source and target message names:

```go
func NewMapper(source, target string) *Mapper
```

### AddField

AddField adds a field mapping to the mapper, resolving the converter from the registry:

```go
func (m *Mapper) AddField(sourceField, targetField string, srcType, dstType TypeInfo, registry *Registry) error
```

Build-time validation occurs here - if no converter is found, the error is returned before code generation.

## Usage Example

```go
import "github.com/jwart212/protoc-gen-go-mapper/internal/graph"

m := graph.NewMapper("User", "User")
r := registry.New()
r.Register(ScalarConverter{})

srcType := types.TypeInfo{Kind: types.KindScalar}
dstType := types.TypeInfo{Kind: types.KindScalar}

err := m.AddField("id", "id", srcType, dstType, r)
if err != nil {
    // Handle error (no converter found)
}
```

## Design Decisions

- **Build-time validation**: Converter resolution happens during graph construction, not during generation.
- **Field order preservation**: Fields are added in the order they appear in the proto definition.
- **No descriptor inspection**: The generator consumes only the Mapper graph, never proto descriptors.

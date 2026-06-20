# converter

Package converter defines the Converter interface used for type conversions.

## Overview

The converter package provides the core interface that all type-specific converters must implement. This is a small, single-purpose interface following Effective Go's guidance.

## Types

### Converter

Converter maps a single field between a proto-side TypeInfo and a db-side TypeInfo:

```go
type Converter interface {
    Match(src, dst TypeInfo) bool
    Priority() int
    Generate(field MappingField) (string, error)
}
```

**Match**: Reports whether this converter handles the (src, dst) type pair.

**Priority**: Breaks ties when multiple converters match the same pair. Higher priority wins. Converters that match narrower, more specific pairs must report higher priority than generic fallbacks.

**Generate**: Emits the Go expression (not a full statement) that performs the conversion for the given field.

### MappingField

MappingField represents a field mapping for code generation:

```go
type MappingField struct {
    SourceField string
    TargetField string
    SourceExpr  string
    TargetExpr  string
}
```

## Usage Example

```go
import "github.com/jwart212/protoc-gen-go-mapper/pkg/converter"

type MyConverter struct{}

func (c MyConverter) Match(src, dst types.TypeInfo) bool {
    return src.Kind == types.KindUUID && dst.Kind == types.KindScalar
}

func (c MyConverter) Priority() int {
    return 10
}

func (c MyConverter) Generate(field converter.MappingField) (string, error) {
    return fmt.Sprintf("uuid.MustParse(%s)", field.SourceExpr), nil
}
```

## Design Decisions

- **Small interface**: Converter has only 3 methods, all focused on one concern: matching and emitting.
- **Stateless**: Implementations must be stateless and safe for concurrent use.
- **Value types**: TypeInfo is passed by value since it is small and immutable.
- **Expression generation**: Generate returns an expression, not a full statement, to allow flexible use in templates.

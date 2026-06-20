# registry

Package registry manages converter registration and priority-based resolution.

## Overview

The registry package implements the converter registry stage of the compiler-style architecture. It manages all registered converters and resolves the best converter for a given type pair using priority-based selection.

## Types

### Registry

Registry manages converter registration and resolution:

```go
type Registry struct {
    converters []Converter
}

func New() *Registry
func (r *Registry) Register(c Converter)
func (r *Registry) Resolve(src, dst TypeInfo) (Converter, error)
```

## Resolution Rules

1. If only one converter matches, it is returned.
2. If multiple converters match, the one with the highest Priority() is returned.
3. If multiple converters match with equal priority, ErrAmbiguousMapping is returned.
4. If no converter matches, ErrNoConverterFound is returned.

## Usage Example

```go
import "github.com/jwart212/protoc-gen-go-mapper/internal/registry"

r := registry.New()
r.Register(ScalarConverter{})
r.Register(UUIDConverter{})

src := types.TypeInfo{Kind: types.KindUUID}
dst := types.TypeInfo{Kind: types.KindScalar}

converter, err := r.Resolve(src, dst)
if err != nil {
    // Handle error
}
```

## Error Handling

- **ErrNoConverterFound**: Returned when no registered converter reports Match == true for the type pair.
- **ErrAmbiguousMapping**: Returned when two or more converters match with equal priority. This is a critical error that prevents nondeterministic behavior.

## Design Decisions

- **Priority-based selection**: Higher priority converters win, allowing specific converters to override generic ones.
- **Fail on ambiguity**: Equal priority ties return ErrAmbiguousMapping rather than silently picking one to prevent nondeterminism.
- **Thread-safe**: The registry may invoke Match concurrently across goroutines during resolution, so converters must be stateless.
- **Slice storage**: Converters are stored in a slice for predictable iteration order during resolution.

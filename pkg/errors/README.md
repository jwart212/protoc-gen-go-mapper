# errors

Package errors defines sentinel errors used throughout the mapper.

## Overview

The errors package provides the canonical error values used across package boundaries. All errors returned across package boundaries must be wrapped with context using fmt.Errorf with %w.

## Sentinel Errors

### ErrNoConverterFound

Returned by the registry when no registered Converter reports Match == true for a type pair.

### ErrAmbiguousMapping

Returned when two or more converters match the same pair with equal priority. This is a critical error that prevents nondeterministic behavior.

### ErrUnsupportedKind

Returned when a Kind is structurally valid but not supported in the current conversion direction (e.g., a NOT NULL DB enum column with no UNSPECIFIED equivalent).

### ErrInvalidConfig

Returned by config validation before parsing begins when the mapper.yaml configuration is invalid.

## Usage Example

```go
import (
    "fmt"
    "github.com/jwart212/protoc-gen-go-mapper/pkg/errors"
)

// Wrap sentinel error with context
err := fmt.Errorf("resolving field %s: %w", fieldName, errors.ErrNoConverterFound)

// Check for specific error using errors.Is
if errors.Is(err, errors.ErrNoConverterFound) {
    // Handle no converter found case
}
```

## Error Contract

1. Every error returned across a package boundary must be wrapped with context using fmt.Errorf with %w.
2. errors.Is / errors.As must work against the sentinel errors.
3. Tests should assert on the sentinel, never on error string content.

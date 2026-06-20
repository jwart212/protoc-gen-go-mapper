package converter

import "gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/types"

// MappingField represents a field mapping for code generation.
type MappingField struct {
	SourceField string
	TargetField string
	SourceExpr  string
	TargetExpr  string
	SourceType  types.TypeInfo
	TargetType  types.TypeInfo
}

// Converter maps a single field between a proto-side TypeInfo and a
// db-side TypeInfo. Implementations must be stateless and safe for
// concurrent use — the registry may invoke Match concurrently across
// goroutines during resolution.
type Converter interface {
	// Match reports whether this converter handles the (src, dst) pair.
	Match(src, dst types.TypeInfo) bool

	// Priority breaks ties when multiple converters Match the same
	// pair. Higher wins. Converters that match a narrower, more
	// specific pair (e.g. NullableConverter wrapping a UUIDConverter)
	// must report a higher priority than the generic fallback.
	Priority() int

	// Generate emits the Go expression (not a full statement) that
	// performs the conversion for the given field.
	Generate(field MappingField) (string, error)
}

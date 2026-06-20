package registry

import (
	"fmt"

	"github.com/jwart212/protoc-gen-go-mapper/pkg/converter"
	"github.com/jwart212/protoc-gen-go-mapper/pkg/types"
)

// EnumConverter handles proto enum ↔ sqlc enum conversions with fallback logic.
type EnumConverter struct{}

// Match returns true for Enum ↔ scalar (string) and Enum ↔ Enum conversions.
func (c EnumConverter) Match(src, dst types.TypeInfo) bool {
	// Enum to string
	if src.Kind == types.KindEnum && dst.Kind == types.KindScalar {
		return true
	}
	// string to Enum
	if src.Kind == types.KindScalar && dst.Kind == types.KindEnum {
		return true
	}
	// Enum to Enum (same enum type)
	if src.Kind == types.KindEnum && dst.Kind == types.KindEnum {
		return true
	}
	return false
}

// Priority returns a higher priority than ScalarConverter for Enum-specific conversions.
func (c EnumConverter) Priority() int {
	return 10
}

// Generate emits the Go expression for Enum conversion with fallback logic.
func (c EnumConverter) Generate(field converter.MappingField) (string, error) {
	// Enum to Enum (direct assignment)
	if field.SourceType.Kind == types.KindEnum && field.TargetType.Kind == types.KindEnum {
		return field.SourceExpr, nil
	}
	// Enum to string (DB → Proto)
	if field.SourceType.Kind == types.KindEnum && field.TargetType.Kind == types.KindScalar {
		return fmt.Sprintf("%s.String()", field.SourceExpr), nil
	}
	// string to Enum (Proto → DB)
	return fmt.Sprintf("%s.String()", field.SourceExpr), nil
}

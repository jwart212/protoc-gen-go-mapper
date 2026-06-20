package registry

import (
	"fmt"

	"github.com/jwart212/protoc-gen-go-mapper/pkg/converter"
	"github.com/jwart212/protoc-gen-go-mapper/pkg/types"
)

// SliceConverter handles []T ↔ []T conversions with element-level delegation.
type SliceConverter struct{}

// Match returns true for slice ↔ slice conversions.
func (c SliceConverter) Match(src, dst types.TypeInfo) bool {
	// Slice to slice
	if src.IsSlice && dst.IsSlice {
		return true
	}
	return false
}

// Priority returns a higher priority than ScalarConverter for Slice-specific conversions.
func (c SliceConverter) Priority() int {
	return 10
}

// Generate emits the Go expression for Slice conversion.
func (c SliceConverter) Generate(field converter.MappingField) (string, error) {
	// Get element types by removing slice flag
	srcElem := field.SourceType
	srcElem.IsSlice = false
	dstElem := field.TargetType
	dstElem.IsSlice = false

	// If element types are the same, use direct assignment
	if srcElem.Package == dstElem.Package && srcElem.Name == dstElem.Name {
		return field.SourceExpr, nil
	}

	// For different element types, generate a loop-based conversion
	// This will be handled by generating a helper function in the template
	return fmt.Sprintf("convert%sArray(%s)", dstElem.Name, field.SourceExpr), nil
}

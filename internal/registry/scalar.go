package registry

import (
	"github.com/jwart212/protoc-gen-go-mapper/pkg/converter"
	"github.com/jwart212/protoc-gen-go-mapper/pkg/types"
)

// ScalarConverter handles basic scalar type conversions (string, int, bool, etc.).
type ScalarConverter struct{}

// Match returns true for scalar types.
func (c ScalarConverter) Match(src, dst types.TypeInfo) bool {
	return src.Kind == types.KindScalar && dst.Kind == types.KindScalar
}

// Priority returns a low priority since this is the generic fallback.
func (c ScalarConverter) Priority() int {
	return 0
}

// Generate emits a simple assignment expression for scalar types.
func (c ScalarConverter) Generate(field converter.MappingField) (string, error) {
	return field.SourceExpr, nil
}

package registry

import (
	"fmt"

	"github.com/jwart212/protoc-gen-go-mapper/pkg/converter"
	"github.com/jwart212/protoc-gen-go-mapper/pkg/types"
)

// DecimalConverter handles decimal.Decimal ↔ string conversions.
type DecimalConverter struct{}

// Match returns true for Decimal ↔ scalar (string) conversions.
func (c DecimalConverter) Match(src, dst types.TypeInfo) bool {
	// Decimal to string
	if src.Kind == types.KindDecimal && dst.Kind == types.KindScalar {
		return true
	}
	// string to Decimal
	if src.Kind == types.KindScalar && dst.Kind == types.KindDecimal {
		return true
	}
	return false
}

// Priority returns a higher priority than ScalarConverter for Decimal-specific conversions.
func (c DecimalConverter) Priority() int {
	return 10
}

// Generate emits the Go expression for Decimal conversion.
func (c DecimalConverter) Generate(field converter.MappingField) (string, error) {
	// Decimal to string
	if field.SourceExpr == "src" && field.TargetExpr == "dst" {
		return fmt.Sprintf("%s.String()", field.SourceExpr), nil
	}
	// string to Decimal
	return fmt.Sprintf("decimal.RequireFromString(%s)", field.SourceExpr), nil
}

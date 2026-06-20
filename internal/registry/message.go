package registry

import (
	"fmt"

	"github.com/jwart212/protoc-gen-go-mapper/pkg/converter"
	"github.com/jwart212/protoc-gen-go-mapper/pkg/types"
)

// MessageConverter handles Message-to-Message conversions (nested messages).
type MessageConverter struct{}

// Match returns true for Message-to-Message conversions.
func (c MessageConverter) Match(src, dst types.TypeInfo) bool {
	return src.Kind == types.KindMessage && dst.Kind == types.KindMessage
}

// Priority returns a low priority for generic message conversions.
func (c MessageConverter) Priority() int {
	return 5
}

// Generate emits a pass-through expression for message types (same type on both sides).
func (c MessageConverter) Generate(field converter.MappingField) (string, error) {
	// For message-to-message conversions with the same type, pass through directly
	if field.SourceType.Name == field.TargetType.Name {
		return field.SourceExpr, nil
	}
	return "", fmt.Errorf("unsupported message conversion: %v -> %v", field.SourceType, field.TargetType)
}

package graph

import (
	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/converter"
	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/types"
)

// NewEdge creates a new field mapping edge.
func NewEdge(sourceField, targetField string, srcType, dstType types.TypeInfo, conv converter.Converter) FieldMapping {
	return FieldMapping{
		SourceField: sourceField,
		TargetField: targetField,
		SourceType:  srcType,
		TargetType:  dstType,
		Converter:   conv,
	}
}

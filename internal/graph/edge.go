package graph

import (
	"github.com/jwart212/protoc-gen-go-mapper/pkg/converter"
	"github.com/jwart212/protoc-gen-go-mapper/pkg/types"
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

package graph

import (
	"fmt"

	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/internal/registry"
	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/converter"
	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/types"
)

// Mapper represents the complete mapping between two message types.
type Mapper struct {
	Source string
	Target string

	Fields []FieldMapping
}

// FieldMapping connects individual fields with their assigned converter.
type FieldMapping struct {
	SourceField string
	TargetField string
	SourceType  types.TypeInfo
	TargetType  types.TypeInfo
	Converter   converter.Converter
}

// NewMapper creates a new Mapper for the given source and target message names.
func NewMapper(source, target string) *Mapper {
	return &Mapper{
		Source: source,
		Target: target,
		Fields: make([]FieldMapping, 0),
	}
}

// AddField adds a field mapping to the mapper.
func (m *Mapper) AddField(sourceField, targetField string, srcType, dstType types.TypeInfo, reg *registry.Registry) error {
	conv, err := reg.Resolve(srcType, dstType)
	if err != nil {
		return fmt.Errorf("adding field %s -> %s: %w", sourceField, targetField, err)
	}

	m.Fields = append(m.Fields, FieldMapping{
		SourceField: sourceField,
		TargetField: targetField,
		SourceType:  srcType,
		TargetType:  dstType,
		Converter:   conv,
	})

	return nil
}

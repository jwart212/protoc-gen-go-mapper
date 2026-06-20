package schema

import "gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/types"

// Field represents a single field in a protobuf message.
type Field struct {
	Name string

	ProtoType types.TypeInfo
	DBType    types.TypeInfo

	// FieldNumber preserves the original proto field declaration order.
	// This is the canonical sort key for deterministic output.
	FieldNumber int32

	Repeated bool
	Optional bool
}

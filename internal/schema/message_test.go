package schema

import (
	"testing"

	"github.com/jwart212/protoc-gen-go-mapper/pkg/types"
)

func TestMessage(t *testing.T) {
	msg := &Message{
		Name: "User",
		Fields: []*Field{
			{
				Name:        "id",
				FieldNumber: 1,
				ProtoType: types.TypeInfo{
					Name: "string",
					Kind: types.KindScalar,
				},
			},
			{
				Name:        "name",
				FieldNumber: 2,
				ProtoType: types.TypeInfo{
					Name: "string",
					Kind: types.KindScalar,
				},
			},
		},
	}

	if msg.Name != "User" {
		t.Errorf("Expected Name to be User, got %s", msg.Name)
	}
	if len(msg.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(msg.Fields))
	}
}

func TestFieldNumberOrdering(t *testing.T) {
	// Test that FieldNumber preserves declaration order
	msg := &Message{
		Name: "User",
		Fields: []*Field{
			{
				Name:        "id",
				FieldNumber: 1,
			},
			{
				Name:        "name",
				FieldNumber: 2,
			},
			{
				Name:        "email",
				FieldNumber: 3,
			},
		},
	}

	// Verify order is preserved
	if msg.Fields[0].FieldNumber != 1 {
		t.Error("Field order not preserved")
	}
	if msg.Fields[1].FieldNumber != 2 {
		t.Error("Field order not preserved")
	}
	if msg.Fields[2].FieldNumber != 3 {
		t.Error("Field order not preserved")
	}
}

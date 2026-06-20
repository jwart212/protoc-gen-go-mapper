package graph

import (
	"testing"

	"github.com/jwart212/protoc-gen-go-mapper/internal/registry"
	"github.com/jwart212/protoc-gen-go-mapper/pkg/types"
)

func TestNewMapper(t *testing.T) {
	m := NewMapper("User", "User")
	if m == nil {
		t.Error("NewMapper() returned nil")
	}
	if m.Source != "User" {
		t.Errorf("Expected Source to be User, got %s", m.Source)
	}
	if m.Target != "User" {
		t.Errorf("Expected Target to be User, got %s", m.Target)
	}
	if m.Fields == nil {
		t.Error("Fields should be initialized")
	}
}

func TestAddField(t *testing.T) {
	m := NewMapper("User", "User")
	r := registry.New()
	r.Register(registry.ScalarConverter{})

	srcType := types.TypeInfo{Kind: types.KindScalar}
	dstType := types.TypeInfo{Kind: types.KindScalar}

	err := m.AddField("id", "id", srcType, dstType, r)
	if err != nil {
		t.Fatalf("AddField() error = %v", err)
	}

	if len(m.Fields) != 1 {
		t.Errorf("Expected 1 field, got %d", len(m.Fields))
	}
	if m.Fields[0].SourceField != "id" {
		t.Errorf("Expected SourceField to be id, got %s", m.Fields[0].SourceField)
	}
}

func TestAddFieldNoConverter(t *testing.T) {
	m := NewMapper("User", "User")
	r := registry.New()

	srcType := types.TypeInfo{Kind: types.KindUUID}
	dstType := types.TypeInfo{Kind: types.KindScalar}

	err := m.AddField("id", "id", srcType, dstType, r)
	if err == nil {
		t.Error("AddField() should return error when no converter found")
	}
}

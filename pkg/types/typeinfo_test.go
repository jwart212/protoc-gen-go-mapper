package types

import "testing"

func TestTypeInfo(t *testing.T) {
	ti := TypeInfo{
		Package:   "github.com/google/uuid",
		Name:      "UUID",
		IsPointer: false,
		IsSlice:   false,
		IsEnum:    false,
		Kind:      KindUUID,
	}

	if ti.Package != "github.com/google/uuid" {
		t.Errorf("Expected Package to be github.com/google/uuid, got %s", ti.Package)
	}
	if ti.Name != "UUID" {
		t.Errorf("Expected Name to be UUID, got %s", ti.Name)
	}
	if ti.Kind != KindUUID {
		t.Errorf("Expected Kind to be KindUUID, got %v", ti.Kind)
	}
}

func TestTypeInfoValueSemantics(t *testing.T) {
	// TypeInfo should be a value type (passed by value)
	original := TypeInfo{
		Package: "pkg",
		Name:    "Type",
		Kind:    KindScalar,
	}

	copied := original
	copied.Name = "Modified"

	if original.Name == "Modified" {
		t.Error("TypeInfo should have value semantics - modification of copy should not affect original")
	}
}

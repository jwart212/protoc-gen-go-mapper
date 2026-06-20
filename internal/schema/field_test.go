package schema

import (
	"testing"

	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/types"
)

func TestField(t *testing.T) {
	field := &Field{
		Name:        "user_id",
		FieldNumber: 1,
		Repeated:    false,
		Optional:    true,
		ProtoType: types.TypeInfo{
			Name: "string",
			Kind: types.KindScalar,
		},
		DBType: types.TypeInfo{
			Name: "string",
			Kind: types.KindScalar,
		},
	}

	if field.Name != "user_id" {
		t.Errorf("Expected Name to be user_id, got %s", field.Name)
	}
	if field.FieldNumber != 1 {
		t.Errorf("Expected FieldNumber to be 1, got %d", field.FieldNumber)
	}
	if field.Repeated {
		t.Error("Expected Repeated to be false")
	}
	if !field.Optional {
		t.Error("Expected Optional to be true")
	}
}

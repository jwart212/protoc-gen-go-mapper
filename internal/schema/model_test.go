package schema

import (
	"testing"
)

func TestModel(t *testing.T) {
	model := &Model{
		Messages: []*Message{
			{
				Name: "User",
				Fields: []*Field{
					{
						Name:        "id",
						FieldNumber: 1,
					},
				},
			},
		},
		Enums: []*Enum{
			{
				Name:   "Status",
				Values: []string{"ACTIVE", "INACTIVE"},
			},
		},
	}

	if len(model.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(model.Messages))
	}
	if len(model.Enums) != 1 {
		t.Errorf("Expected 1 enum, got %d", len(model.Enums))
	}
}

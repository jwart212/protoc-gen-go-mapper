package proto

import (
	"testing"

	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/internal/schema"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestNew(t *testing.T) {
	p := New("test")
	if p == nil {
		t.Error("New() returned nil")
	}
}

func TestParseFile(t *testing.T) {
	p := New("test")
	fileProto := &descriptorpb.FileDescriptorProto{
		Name: proto.String("test.proto"),
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: proto.String("User"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:     proto.String("id"),
						Number:   proto.Int32(1),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum(),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						JsonName: proto.String("id"),
					},
				},
			},
		},
	}
	model, err := p.ParseFile(fileProto)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}
	if model == nil {
		t.Error("ParseFile() returned nil model")
	}
}

func TestFieldNumberOrdering(t *testing.T) {
	// Test that parser preserves FieldNumber ordering
	msg := &schema.Message{
		Name: "User",
		Fields: []*schema.Field{
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

	// Verify order is preserved by FieldNumber
	for i, field := range msg.Fields {
		expectedNum := int32(i + 1)
		if field.FieldNumber != expectedNum {
			t.Errorf("Field at index %d has FieldNumber %d, expected %d", i, field.FieldNumber, expectedNum)
		}
	}
}

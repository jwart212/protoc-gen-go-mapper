package schema_test

import (
	"fmt"

	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/internal/schema"
	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/types"
)

func ExampleMessage() {
	msg := &schema.Message{
		Name: "User",
		Fields: []*schema.Field{
			{
				Name:        "id",
				FieldNumber: 1,
				ProtoType:   types.TypeInfo{Name: "string"},
			},
		},
	}
	fmt.Printf("Message: %s with %d fields", msg.Name, len(msg.Fields))
	// Output: Message: User with 1 fields
}

func ExampleField() {
	field := &schema.Field{
		Name:        "user_id",
		FieldNumber: 1,
	}
	fmt.Printf("Field: %s (number: %d)", field.Name, field.FieldNumber)
	// Output: Field: user_id (number: 1)
}

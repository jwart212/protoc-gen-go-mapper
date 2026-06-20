package generator_test

import (
	"fmt"

	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/internal/generator"
	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/internal/graph"
	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/internal/schema"
)

func ExampleNew() {
	g := generator.New()
	fmt.Printf("Generator created: %v", g != nil)
	// Output: Generator created: true
}

func ExampleGenerator_Generate() {
	g := generator.New()

	msg := &schema.Message{
		Name: "User",
		Fields: []*schema.Field{
			{Name: "id", FieldNumber: 1},
		},
	}

	protoToDB := graph.NewMapper("User", "User")
	dbToProto := graph.NewMapper("User", "User")
	typeMappings := map[string]string{}

	code, _ := g.Generate(msg, protoToDB, dbToProto, typeMappings)
	fmt.Printf("Generated %d characters of code", len(code))
	// Output: Generated 150 characters of code
}

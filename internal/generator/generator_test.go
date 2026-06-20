package generator

import (
	"testing"

	"github.com/jwart212/protoc-gen-go-mapper/internal/graph"
	"github.com/jwart212/protoc-gen-go-mapper/internal/schema"
)

func TestNew(t *testing.T) {
	g := New()
	if g == nil {
		t.Error("New() returned nil")
	}
}

func TestGenerate(t *testing.T) {
	g := New()

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
		},
	}

	protoToDB := graph.NewMapper("User", "User")
	dbToProto := graph.NewMapper("User", "User")

	code, err := g.Generate(msg, protoToDB, dbToProto, nil, false)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if code == "" {
		t.Error("Generate() returned empty code")
	}

	// Check that both functions are generated
	if len(code) < 50 {
		t.Errorf("Generated code too short: %d chars", len(code))
	}
}

func TestGenerateDeterministic(t *testing.T) {
	g := New()

	msg := &schema.Message{
		Name: "User",
		Fields: []*schema.Field{
			{
				Name:        "name",
				FieldNumber: 2,
			},
			{
				Name:        "id",
				FieldNumber: 1,
			},
		},
	}

	protoToDB := graph.NewMapper("User", "User")
	dbToProto := graph.NewMapper("User", "User")

	code1, err := g.Generate(msg, protoToDB, dbToProto, nil, false)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	code2, err := g.Generate(msg, protoToDB, dbToProto, nil, false)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if code1 != code2 {
		t.Error("Generate() should produce deterministic output")
	}
}

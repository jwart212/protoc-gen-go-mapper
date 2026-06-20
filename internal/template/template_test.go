package template

import (
	"testing"
)

func TestNew(t *testing.T) {
	tmpl := New()
	if tmpl == nil {
		t.Error("New() returned nil")
	}
}

func TestLoad(t *testing.T) {
	tmpl := New()

	err := tmpl.Load("test", "Hello {{.Name}}")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
}

func TestExecute(t *testing.T) {
	tmpl := New()

	err := tmpl.Load("test", "Hello {{.Name}}")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	type Data struct {
		Name string
	}

	result, err := tmpl.Execute("test", Data{Name: "World"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result != "Hello World" {
		t.Errorf("Execute() = %v, want Hello World", result)
	}
}

package template_test

import (
	"fmt"

	"github.com/jwart212/protoc-gen-go-mapper/internal/template"
)

func ExampleNew() {
	tmpl := template.New()
	fmt.Printf("Template created: %v", tmpl != nil)
	// Output: Template created: true
}

func ExampleTemplate_Load() {
	tmpl := template.New()
	tmpl.Load("test", "Hello {{.Name}}")
	fmt.Println("Template loaded")
	// Output: Template loaded
}

func ExampleTemplate_Execute() {
	tmpl := template.New()
	tmpl.Load("test", "Hello {{.Name}}")

	result, _ := tmpl.Execute("test", struct{ Name string }{Name: "World"})
	fmt.Println(result)
	// Output: Hello World
}

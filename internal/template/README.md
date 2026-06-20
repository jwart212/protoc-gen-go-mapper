# template

Package template provides a simple wrapper around Go's text/template for code generation.

## Overview

The template package wraps Go's standard text/template library to provide a clean interface for loading and executing code generation templates.

## Types

### Template

Template manages Go code generation templates:

```go
type Template struct {
    templates map[string]*template.Template
}

func New() *Template
func (t *Template) Load(name, content string) error
func (t *Template) Execute(name string, data interface{}) (string, error)
```

## Usage Example

```go
import "gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/internal/template"

tmpl := template.New()
tmpl.Load("mapper", "func To{{.Name}}(src {{.Type}}) {{.Type}} { return src }")

result, err := tmpl.Execute("mapper", struct {
    Name string
    Type string
}{Name: "User", Type: "string"})
```

## Design Decisions

- **Simple wrapper**: Thin wrapper around text/template with no additional logic.
- **No business logic**: Templates contain only presentation logic, no business logic.
- **Cached templates**: Templates are parsed once and cached for reuse.

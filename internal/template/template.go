package template

import (
	"bytes"
	"text/template"
)

// Template manages Go code generation templates.
type Template struct {
	templates map[string]*template.Template
}

// New creates a new Template instance.
func New() *Template {
	return &Template{
		templates: make(map[string]*template.Template),
	}
}

// Load loads a template with the given name and content.
func (t *Template) Load(name, content string) error {
	tmpl, err := template.New(name).Parse(content)
	if err != nil {
		return err
	}
	t.templates[name] = tmpl
	return nil
}

// Execute executes the template with the given data.
func (t *Template) Execute(name string, data interface{}) (string, error) {
	tmpl, ok := t.templates[name]
	if !ok {
		return "", nil
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

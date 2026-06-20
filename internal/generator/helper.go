package generator

import (
	"fmt"
	"strings"
)

// Helper provides utility functions for code generation.
type Helper struct{}

// NewHelper creates a new Helper instance.
func NewHelper() *Helper {
	return &Helper{}
}

// ToCamelCase converts snake_case to camelCase.
func (h *Helper) ToCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
		}
	}
	return strings.Join(parts, "")
}

// ToPascalCase converts snake_case to PascalCase.
func (h *Helper) ToPascalCase(s string) string {
	camel := h.ToCamelCase(s)
	if len(camel) > 0 {
		return strings.ToUpper(camel[:1]) + camel[1:]
	}
	return camel
}

// MapSlice generates code for mapping a slice.
func (h *Helper) MapSlice(src, elemType string) string {
	return fmt.Sprintf("MapSlice(%s, func(item %s) %s { return item })", src, elemType, elemType)
}

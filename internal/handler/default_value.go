package handler

import (
	"fmt"
	"github.com/jwart212/protoc-gen-go-mapper/internal/schema"
	"strings"
)

// DefaultValueHandler sets default values for fields that are not present in the DB struct.
// For example, setting Children field to []*CategoryTreeNode{}.
type DefaultValueHandler struct {
	fieldName    string // Field name to match (case-insensitive)
	dbTypeNames  []string // DB type names to match
	defaultValue string // The default value expression (e.g., "[]*CategoryTreeNode{}")
}

// NewDefaultValueHandler creates a new DefaultValueHandler.
func NewDefaultValueHandler(fieldName string, dbTypeNames []string, defaultValue string) *DefaultValueHandler {
	return &DefaultValueHandler{
		fieldName:    fieldName,
		dbTypeNames:  dbTypeNames,
		defaultValue: defaultValue,
	}
}

// Match returns true if the field name and DB type match the configuration.
func (h *DefaultValueHandler) Match(field *schema.Field, dbTypeName string) bool {
	if h.fieldName == "" {
		return false
	}
	if strings.ToLower(field.Name) != strings.ToLower(h.fieldName) {
		return false
	}
	if len(h.dbTypeNames) == 0 {
		return true
	}
	for _, typeName := range h.dbTypeNames {
		if dbTypeName == typeName {
			return true
		}
	}
	return false
}

// GenerateToProto generates the default value assignment for DB -> Proto conversion.
// Example: "Children: []*CategoryTreeNode{}"
func (h *DefaultValueHandler) GenerateToProto(field *schema.Field, dbFieldName, protoFieldName, dbTypeName string) (string, error) {
	if h.defaultValue == "" {
		return "", fmt.Errorf("defaultValue not configured for DefaultValueHandler")
	}
	return fmt.Sprintf("%s: %s", protoFieldName, h.defaultValue), nil
}

// GenerateToDB returns empty string since default values are typically only for Proto output.
func (h *DefaultValueHandler) GenerateToDB(field *schema.Field, dbFieldName, protoFieldName, dbTypeName string) (string, error) {
	return "", nil
}

// Priority returns a medium-high priority for default value handlers.
func (h *DefaultValueHandler) Priority() int {
	return 75
}

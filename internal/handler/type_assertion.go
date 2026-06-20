package handler

import (
	"fmt"
	"strings"

	"github.com/jwart212/protoc-gen-go-mapper/internal/schema"
)

// TypeAssertionHandler handles type assertions for fields that need explicit type conversion.
// For example, converting interface{} to []string.
type TypeAssertionHandler struct {
	fieldName   string   // Field name to match (case-insensitive)
	dbTypeNames []string // DB type names to match
	assertType  string   // The type to assert to (e.g., "[]string")
}

// NewTypeAssertionHandler creates a new TypeAssertionHandler.
func NewTypeAssertionHandler(fieldName string, dbTypeNames []string, assertType string) *TypeAssertionHandler {
	return &TypeAssertionHandler{
		fieldName:   fieldName,
		dbTypeNames: dbTypeNames,
		assertType:  assertType,
	}
}

// Match returns true if the field name and DB type match the configuration.
func (h *TypeAssertionHandler) Match(field *schema.Field, dbTypeName string) bool {
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

// GenerateToProto generates the type assertion code for DB -> Proto conversion.
// Returns just the field name with type assertion, without variable prefix.
// Example: "Path: .Path.([]string)" - caller should prepend variable name
// Note: This generates a type assertion that will panic at runtime if the type doesn't match.
// Users should ensure the DB field type matches the assertType configuration.
func (h *TypeAssertionHandler) GenerateToProto(field *schema.Field, dbFieldName, protoFieldName, dbTypeName string) (string, error) {
	if h.assertType == "" {
		return "", fmt.Errorf("assertType not configured for TypeAssertionHandler")
	}
	return fmt.Sprintf("%s: .%s.(%s)", protoFieldName, dbFieldName, h.assertType), nil
}

// GenerateToDB generates the direct assignment code for Proto -> DB conversion.
// Type assertions are typically only needed for DB -> Proto direction.
func (h *TypeAssertionHandler) GenerateToDB(field *schema.Field, dbFieldName, protoFieldName, dbTypeName string) (string, error) {
	return fmt.Sprintf("%s: src.%s", dbFieldName, protoFieldName), nil
}

// Priority returns a medium priority for type assertion handlers.
func (h *TypeAssertionHandler) Priority() int {
	return 50
}

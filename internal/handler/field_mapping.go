package handler

import (
	"fmt"
	"strings"

	"github.com/jwart212/protoc-gen-go-mapper/internal/schema"
)

// FieldMappingHandler maps one field to another field in the source/target.
// For example, mapping ParentId to use TenantID from the source.
type FieldMappingHandler struct {
	fieldName   string   // Field name to match (case-insensitive)
	dbTypeNames []string // DB type names to match
	messageName string   // Message name to match (case-insensitive)
	toProtoExpr string   // Custom expression for DB -> Proto conversion
	toDBExpr    string   // Custom expression for Proto -> DB conversion
}

// NewFieldMappingHandler creates a new FieldMappingHandler.
func NewFieldMappingHandler(fieldName, messageName string, dbTypeNames []string, toProtoExpr, toDBExpr string) *FieldMappingHandler {
	return &FieldMappingHandler{
		fieldName:   fieldName,
		dbTypeNames: dbTypeNames,
		messageName: messageName,
		toProtoExpr: toProtoExpr,
		toDBExpr:    toDBExpr,
	}
}

// Match returns true if the field name and DB type match the configuration.
// Message name filtering should be handled at a higher level (in the generator).
func (h *FieldMappingHandler) Match(field *schema.Field, dbTypeName string) bool {
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

// GenerateToProto generates the custom expression for DB -> Proto conversion.
// The expression can use placeholders like {DbField} and {ProtoField}.
func (h *FieldMappingHandler) GenerateToProto(field *schema.Field, dbFieldName, protoFieldName, dbTypeName string) (string, error) {
	if h.toProtoExpr == "" {
		return "", fmt.Errorf("toProtoExpr not configured for FieldMappingHandler")
	}
	expr := strings.ReplaceAll(h.toProtoExpr, "{DbField}", dbFieldName)
	expr = strings.ReplaceAll(expr, "{ProtoField}", protoFieldName)
	return expr, nil
}

// GenerateToDB generates the custom expression for Proto -> DB conversion.
// The expression can use placeholders like {DbField} and {ProtoField}.
func (h *FieldMappingHandler) GenerateToDB(field *schema.Field, dbFieldName, protoFieldName, dbTypeName string) (string, error) {
	if h.toDBExpr == "" {
		return "", fmt.Errorf("toDBExpr not configured for FieldMappingHandler")
	}
	expr := strings.ReplaceAll(h.toDBExpr, "{DbField}", dbFieldName)
	expr = strings.ReplaceAll(expr, "{ProtoField}", protoFieldName)
	return expr, nil
}

// Priority returns a high priority for field mapping handlers.
func (h *FieldMappingHandler) Priority() int {
	return 80
}

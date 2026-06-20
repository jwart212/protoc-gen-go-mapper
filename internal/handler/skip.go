package handler

import (
	"github.com/jwart212/protoc-gen-go-mapper/internal/schema"
	"strings"
)

// SkipHandler skips fields that are not present in the DB struct.
// It can be configured to skip specific field names (case-insensitive).
type SkipHandler struct {
	fieldNames []string // Field names to skip (case-insensitive)
}

// NewSkipHandler creates a new SkipHandler with the given field names to skip.
func NewSkipHandler(fieldNames ...string) *SkipHandler {
	return &SkipHandler{
		fieldNames: fieldNames,
	}
}

// Match returns true if the field name matches any of the configured field names (case-insensitive).
func (h *SkipHandler) Match(field *schema.Field, dbTypeName string) bool {
	if len(h.fieldNames) == 0 {
		return false
	}
	fieldNameLower := strings.ToLower(field.Name)
	for _, name := range h.fieldNames {
		if strings.ToLower(name) == fieldNameLower {
			return true
		}
	}
	return false
}

// GenerateToProto returns an empty string since the field is skipped.
// The caller should not generate any code for this field.
func (h *SkipHandler) GenerateToProto(field *schema.Field, dbFieldName, protoFieldName, dbTypeName string) (string, error) {
	return "", nil
}

// GenerateToDB returns an empty string since the field is skipped.
// The caller should not generate any code for this field.
func (h *SkipHandler) GenerateToDB(field *schema.Field, dbFieldName, protoFieldName, dbTypeName string) (string, error) {
	return "", nil
}

// Priority returns a high priority to ensure skip handlers are checked first.
func (h *SkipHandler) Priority() int {
	return 100
}

package handler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jwart212/protoc-gen-go-mapper/internal/schema"
)

// TypeConversionHandler handles type-based field conversions.
// It matches fields based on their proto type, DB type, optionality, and patterns.
type TypeConversionHandler struct {
	protoType           string
	dbType              string
	isOptional          bool
	toProtoExpr         string
	toDBExpr            string
	priority            int
	matchFieldPattern   *regexp.Regexp
	matchMessagePattern *regexp.Regexp
	pointerStrategy     string
}

// NewTypeConversionHandler creates a new TypeConversionHandler.
func NewTypeConversionHandler(protoType, dbType string, isOptional bool, toProtoExpr, toDBExpr string, priority int) *TypeConversionHandler {
	return &TypeConversionHandler{
		protoType:   protoType,
		dbType:      dbType,
		isOptional:  isOptional,
		toProtoExpr: toProtoExpr,
		toDBExpr:    toDBExpr,
		priority:    priority,
	}
}

// NewTypeConversionHandlerWithPatterns creates a new TypeConversionHandler with pattern matching.
func NewTypeConversionHandlerWithPatterns(protoType, dbType string, isOptional bool, toProtoExpr, toDBExpr string, priority int, matchFieldPattern, matchMessagePattern, pointerStrategy string) (*TypeConversionHandler, error) {
	h := &TypeConversionHandler{
		protoType:       protoType,
		dbType:          dbType,
		isOptional:      isOptional,
		toProtoExpr:     toProtoExpr,
		toDBExpr:        toDBExpr,
		priority:        priority,
		pointerStrategy: pointerStrategy,
	}

	if matchFieldPattern != "" {
		// Make pattern case-insensitive by adding (?i) flag
		re, err := regexp.Compile("(?i)" + matchFieldPattern)
		if err != nil {
			return nil, fmt.Errorf("invalid field pattern %q: %w", matchFieldPattern, err)
		}
		h.matchFieldPattern = re
	}

	if matchMessagePattern != "" {
		re, err := regexp.Compile("(?i)" + matchMessagePattern)
		if err != nil {
			return nil, fmt.Errorf("invalid message pattern %q: %w", matchMessagePattern, err)
		}
		h.matchMessagePattern = re
	}

	return h, nil
}

// Match returns true if the field's type and optionality match the configuration.
func (h *TypeConversionHandler) Match(field *schema.Field, dbTypeName string) bool {
	// Check field pattern first (highest priority)
	if h.matchFieldPattern != nil {
		if !h.matchFieldPattern.MatchString(field.Name) {
			return false
		}
	}

	// Check proto type - handle both simple names and full package names
	protoTypeName := field.ProtoType.Name
	if h.protoType != "" {
		// Handle "google.protobuf.Timestamp" vs "Timestamp"
		if h.protoType == "google.protobuf.Timestamp" {
			if protoTypeName != "Timestamp" && protoTypeName != "google.protobuf.Timestamp" {
				return false
			}
		} else if protoTypeName != h.protoType {
			return false
		}
	}

	// Check DB type
	if h.dbType != "" && field.DBType.Name != h.dbType {
		return false
	}

	// Check optionality
	if h.isOptional != field.Optional {
		return false
	}

	return true
}

// GenerateToProto generates the conversion code for DB -> Proto direction.
// Supports placeholders: {dbField}, {protoField}, {variable}
func (h *TypeConversionHandler) GenerateToProto(field *schema.Field, dbFieldName, protoFieldName, dbTypeName string) (string, error) {
	if h.toProtoExpr == "" {
		return "", fmt.Errorf("toProtoExpr not configured for TypeConversionHandler")
	}
	expr := strings.ReplaceAll(h.toProtoExpr, "{dbField}", dbFieldName)
	expr = strings.ReplaceAll(expr, "{protoField}", protoFieldName)
	return fmt.Sprintf("%s: %s", protoFieldName, expr), nil
}

// GenerateToDB generates the conversion code for Proto -> DB direction.
// Supports placeholders: {dbField}, {protoField}, {variable}
func (h *TypeConversionHandler) GenerateToDB(field *schema.Field, dbFieldName, protoFieldName, dbTypeName string) (string, error) {
	if h.toDBExpr == "" {
		return "", fmt.Errorf("toDBExpr not configured for TypeConversionHandler")
	}
	expr := strings.ReplaceAll(h.toDBExpr, "{dbField}", dbFieldName)
	expr = strings.ReplaceAll(expr, "{protoField}", protoFieldName)
	return fmt.Sprintf("%s: %s", dbFieldName, expr), nil
}

// Priority returns the handler's priority.
func (h *TypeConversionHandler) Priority() int {
	return h.priority
}

// PointerStrategy returns the pointer handling strategy.
func (h *TypeConversionHandler) PointerStrategy() string {
	return h.pointerStrategy
}

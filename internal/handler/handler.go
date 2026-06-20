package handler

import (
	"github.com/jwart212/protoc-gen-go-mapper/internal/schema"
)

// FieldHandler defines the interface for custom field handling during code generation.
// Implementations can provide specialized logic for specific fields or types.
type FieldHandler interface {
	// Match returns true if this handler should process the given field.
	// The dbTypeName is provided for context about the target DB type.
	Match(field *schema.Field, dbTypeName string) bool

	// GenerateToProto generates the Go code for DB -> Proto field conversion.
	// Returns the code expression (not a full statement) for the field assignment.
	GenerateToProto(field *schema.Field, dbFieldName, protoFieldName, dbTypeName string) (string, error)

	// GenerateToDB generates the Go code for Proto -> DB field conversion.
	// Returns the code expression (not a full statement) for the field assignment.
	GenerateToDB(field *schema.Field, dbFieldName, protoFieldName, dbTypeName string) (string, error)

	// Priority determines handler precedence when multiple handlers match.
	// Higher priority handlers are checked first. Return 0 for default priority.
	Priority() int
}

// HandlerRegistry manages a collection of field handlers.
type HandlerRegistry struct {
	handlers []FieldHandler
}

// NewHandlerRegistry creates a new empty handler registry.
func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make([]FieldHandler, 0),
	}
}

// Register adds a handler to the registry.
func (r *HandlerRegistry) Register(handler FieldHandler) {
	r.handlers = append(r.handlers, handler)
}

// Find finds the best matching handler for the given field.
// Returns nil if no handler matches.
func (r *HandlerRegistry) Find(field *schema.Field, dbTypeName string) FieldHandler {
	var best FieldHandler
	var bestPriority int
	found := false

	for _, handler := range r.handlers {
		if handler.Match(field, dbTypeName) {
			priority := handler.Priority()
			if !found {
				best = handler
				bestPriority = priority
				found = true
			} else {
				if priority > bestPriority {
					best = handler
					bestPriority = priority
				}
			}
		}
	}

	return best
}

package registry

import (
	"fmt"

	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/converter"
	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/errors"
	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/types"
)

// Registry manages converter registration and resolution.
type Registry struct {
	converters []converter.Converter
}

// New creates a new Registry instance.
func New() *Registry {
	return &Registry{
		converters: make([]converter.Converter, 0),
	}
}

// Register adds a converter to the registry.
func (r *Registry) Register(c converter.Converter) {
	r.converters = append(r.converters, c)
}

// Resolve finds the best converter for the given type pair.
// Returns ErrNoConverterFound if no converter matches.
// Returns ErrAmbiguousMapping if multiple converters match with equal priority.
func (r *Registry) Resolve(src, dst types.TypeInfo) (converter.Converter, error) {
	var best converter.Converter
	var bestPriority int
	found := false

	for _, c := range r.converters {
		if c.Match(src, dst) {
			priority := c.Priority()
			if !found {
				best = c
				bestPriority = priority
				found = true
			} else {
				if priority == bestPriority {
					return nil, fmt.Errorf("resolving %v -> %v: %w", src, dst, errors.ErrAmbiguousMapping)
				}
				if priority > bestPriority {
					best = c
					bestPriority = priority
				}
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("resolving %v -> %v: %w", src, dst, errors.ErrNoConverterFound)
	}

	return best, nil
}

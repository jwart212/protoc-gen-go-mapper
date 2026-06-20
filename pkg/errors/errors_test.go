package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestSentinelErrors(t *testing.T) {
	tests := []struct {
		name  string
		err   error
		check func(error) bool
	}{
		{
			name:  "ErrNoConverterFound",
			err:   ErrNoConverterFound,
			check: func(e error) bool { return errors.Is(e, ErrNoConverterFound) },
		},
		{
			name:  "ErrAmbiguousMapping",
			err:   ErrAmbiguousMapping,
			check: func(e error) bool { return errors.Is(e, ErrAmbiguousMapping) },
		},
		{
			name:  "ErrUnsupportedKind",
			err:   ErrUnsupportedKind,
			check: func(e error) bool { return errors.Is(e, ErrUnsupportedKind) },
		},
		{
			name:  "ErrInvalidConfig",
			err:   ErrInvalidConfig,
			check: func(e error) bool { return errors.Is(e, ErrInvalidConfig) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.check(tt.err) {
				t.Errorf("errors.Is check failed for %s", tt.name)
			}

			// Test wrapping with context
			wrapped := fmt.Errorf("context: %w", tt.err)
			if !tt.check(wrapped) {
				t.Errorf("errors.Is check failed for wrapped %s", tt.name)
			}
		})
	}
}

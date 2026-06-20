package errors

import "errors"

var (
	// ErrNoConverterFound is returned by the registry when no
	// registered Converter reports Match == true for a type pair.
	ErrNoConverterFound = errors.New("mapper: no converter found for type pair")

	// ErrAmbiguousMapping is returned when two or more converters
	// match the same pair with equal priority.
	ErrAmbiguousMapping = errors.New("mapper: multiple converters matched with equal priority")

	// ErrUnsupportedKind is returned when a Kind is structurally
	// valid but not supported in the current conversion direction
	// (e.g. a NOT NULL DB enum column with no UNSPECIFIED equivalent).
	ErrUnsupportedKind = errors.New("mapper: unsupported kind for this conversion direction")

	// ErrInvalidConfig is returned by config validation before
	// parsing begins.
	ErrInvalidConfig = errors.New("mapper: invalid mapper.yaml configuration")
)

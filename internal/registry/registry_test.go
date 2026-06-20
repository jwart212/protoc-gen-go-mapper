package registry

import (
	"errors"
	"testing"

	"github.com/jwart212/protoc-gen-go-mapper/pkg/converter"
	mappererrors "github.com/jwart212/protoc-gen-go-mapper/pkg/errors"
	"github.com/jwart212/protoc-gen-go-mapper/pkg/types"
)

// mockConverter is a test helper that implements the Converter interface
type mockConverter struct {
	matchFunc    func(src, dst types.TypeInfo) bool
	priorityFunc func() int
	generateFunc func(field converter.MappingField) (string, error)
}

func (m *mockConverter) Match(src, dst types.TypeInfo) bool {
	return m.matchFunc(src, dst)
}

func (m *mockConverter) Priority() int {
	return m.priorityFunc()
}

func (m *mockConverter) Generate(field converter.MappingField) (string, error) {
	return m.generateFunc(field)
}

func TestNew(t *testing.T) {
	r := New()
	if r == nil {
		t.Error("New() returned nil")
	}
}

func TestRegister(t *testing.T) {
	r := New()
	c := ScalarConverter{}
	r.Register(c)
	if len(r.converters) != 1 {
		t.Errorf("Expected 1 converter, got %d", len(r.converters))
	}
}

func TestResolve(t *testing.T) {
	r := New()
	r.Register(ScalarConverter{})

	src := types.TypeInfo{Kind: types.KindScalar}
	dst := types.TypeInfo{Kind: types.KindScalar}

	c, err := r.Resolve(src, dst)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if c == nil {
		t.Error("Resolve() returned nil converter")
	}
}

func TestResolveNoConverter(t *testing.T) {
	r := New()

	src := types.TypeInfo{Kind: types.KindUUID}
	dst := types.TypeInfo{Kind: types.KindScalar}

	_, err := r.Resolve(src, dst)
	if err == nil {
		t.Error("Resolve() should return error when no converter matches")
	}
	if !errors.Is(err, mappererrors.ErrNoConverterFound) {
		t.Errorf("Resolve() error should wrap ErrNoConverterFound, got %v", err)
	}
}

func TestResolveAmbiguous(t *testing.T) {
	// Create a mock converter with same priority as ScalarConverter
	ambConv := &mockConverter{
		matchFunc: func(src, dst types.TypeInfo) bool {
			return src.Kind == types.KindScalar && dst.Kind == types.KindScalar
		},
		priorityFunc: func() int {
			return 0
		},
		generateFunc: func(field converter.MappingField) (string, error) {
			return field.SourceExpr, nil
		},
	}

	r := New()
	r.Register(ScalarConverter{})
	r.Register(ambConv)

	src := types.TypeInfo{Kind: types.KindScalar}
	dst := types.TypeInfo{Kind: types.KindScalar}

	_, err := r.Resolve(src, dst)
	if err == nil {
		t.Error("Resolve() should return error for ambiguous priority")
	}
	if !errors.Is(err, mappererrors.ErrAmbiguousMapping) {
		t.Errorf("Resolve() error should wrap ErrAmbiguousMapping, got %v", err)
	}
}

func TestResolvePriority(t *testing.T) {
	// Create high-priority converter
	highConv := &mockConverter{
		matchFunc: func(src, dst types.TypeInfo) bool {
			return src.Kind == types.KindScalar && dst.Kind == types.KindScalar
		},
		priorityFunc: func() int {
			return 10
		},
		generateFunc: func(field converter.MappingField) (string, error) {
			return "high_priority", nil
		},
	}

	r := New()
	r.Register(ScalarConverter{})
	r.Register(highConv)

	src := types.TypeInfo{Kind: types.KindScalar}
	dst := types.TypeInfo{Kind: types.KindScalar}

	c, err := r.Resolve(src, dst)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if c.Priority() != 10 {
		t.Errorf("Expected high priority converter (10), got %d", c.Priority())
	}
}

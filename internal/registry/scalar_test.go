package registry

import (
	"testing"

	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/converter"
	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/types"
)

func TestScalarConverterMatch(t *testing.T) {
	c := ScalarConverter{}

	tests := []struct {
		name string
		src  types.TypeInfo
		dst  types.TypeInfo
		want bool
	}{
		{
			name: "scalar to scalar",
			src:  types.TypeInfo{Kind: types.KindScalar},
			dst:  types.TypeInfo{Kind: types.KindScalar},
			want: true,
		},
		{
			name: "uuid to scalar",
			src:  types.TypeInfo{Kind: types.KindUUID},
			dst:  types.TypeInfo{Kind: types.KindScalar},
			want: false,
		},
		{
			name: "scalar to uuid",
			src:  types.TypeInfo{Kind: types.KindScalar},
			dst:  types.TypeInfo{Kind: types.KindUUID},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Match(tt.src, tt.dst); got != tt.want {
				t.Errorf("ScalarConverter.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScalarConverterPriority(t *testing.T) {
	c := ScalarConverter{}
	if c.Priority() != 0 {
		t.Errorf("Expected priority 0, got %d", c.Priority())
	}
}

func TestScalarConverterGenerate(t *testing.T) {
	c := ScalarConverter{}
	field := converter.MappingField{
		SourceExpr: "src.Field",
	}

	got, err := c.Generate(field)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got != "src.Field" {
		t.Errorf("Generate() = %v, want %v", got, "src.Field")
	}
}

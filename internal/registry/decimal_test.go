package registry

import (
	"testing"

	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/converter"
	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/types"
)

func TestDecimalConverterMatch(t *testing.T) {
	c := DecimalConverter{}

	tests := []struct {
		name string
		src  types.TypeInfo
		dst  types.TypeInfo
		want bool
	}{
		{
			name: "Decimal to string",
			src:  types.TypeInfo{Kind: types.KindDecimal},
			dst:  types.TypeInfo{Kind: types.KindScalar},
			want: true,
		},
		{
			name: "string to Decimal",
			src:  types.TypeInfo{Kind: types.KindScalar},
			dst:  types.TypeInfo{Kind: types.KindDecimal},
			want: true,
		},
		{
			name: "Decimal to Decimal",
			src:  types.TypeInfo{Kind: types.KindDecimal},
			dst:  types.TypeInfo{Kind: types.KindDecimal},
			want: false,
		},
		{
			name: "scalar to scalar",
			src:  types.TypeInfo{Kind: types.KindScalar},
			dst:  types.TypeInfo{Kind: types.KindScalar},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Match(tt.src, tt.dst); got != tt.want {
				t.Errorf("DecimalConverter.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimalConverterPriority(t *testing.T) {
	c := DecimalConverter{}
	if c.Priority() != 10 {
		t.Errorf("Expected priority 10, got %d", c.Priority())
	}
}

func TestDecimalConverterGenerate(t *testing.T) {
	c := DecimalConverter{}

	field := converter.MappingField{
		SourceExpr: "src",
		TargetExpr: "dst",
	}

	got, err := c.Generate(field)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got == "" {
		t.Error("Generate() returned empty string")
	}
}

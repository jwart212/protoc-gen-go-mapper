package registry

import (
	"testing"

	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/converter"
	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/types"
)

func TestEnumConverterMatch(t *testing.T) {
	c := EnumConverter{}

	tests := []struct {
		name string
		src  types.TypeInfo
		dst  types.TypeInfo
		want bool
	}{
		{
			name: "Enum to string",
			src:  types.TypeInfo{Kind: types.KindEnum},
			dst:  types.TypeInfo{Kind: types.KindScalar},
			want: true,
		},
		{
			name: "string to Enum",
			src:  types.TypeInfo{Kind: types.KindScalar},
			dst:  types.TypeInfo{Kind: types.KindEnum},
			want: true,
		},
		{
			name: "Enum to Enum",
			src:  types.TypeInfo{Kind: types.KindEnum},
			dst:  types.TypeInfo{Kind: types.KindEnum},
			want: true,
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
				t.Errorf("EnumConverter.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnumConverterPriority(t *testing.T) {
	c := EnumConverter{}
	if c.Priority() != 10 {
		t.Errorf("Expected priority 10, got %d", c.Priority())
	}
}

func TestEnumConverterGenerate(t *testing.T) {
	c := EnumConverter{}

	tests := []struct {
		name  string
		field converter.MappingField
		want  string
	}{
		{
			name: "Enum to string",
			field: converter.MappingField{
				SourceExpr: "src",
				TargetExpr: "dst",
				SourceType: types.TypeInfo{Kind: types.KindEnum},
				TargetType: types.TypeInfo{Kind: types.KindScalar},
			},
			want: "src.String()",
		},
		{
			name: "string to Enum",
			field: converter.MappingField{
				SourceExpr: "src",
				TargetExpr: "dst",
				SourceType: types.TypeInfo{Kind: types.KindScalar},
				TargetType: types.TypeInfo{Kind: types.KindEnum},
			},
			want: "src.String()",
		},
		{
			name: "Enum to Enum",
			field: converter.MappingField{
				SourceExpr: "src",
				TargetExpr: "dst",
				SourceType: types.TypeInfo{Kind: types.KindEnum},
				TargetType: types.TypeInfo{Kind: types.KindEnum},
			},
			want: "src",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Generate(tt.field)
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}

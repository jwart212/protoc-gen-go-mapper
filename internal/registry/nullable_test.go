package registry

import (
	"testing"

	"github.com/jwart212/protoc-gen-go-mapper/pkg/converter"
	"github.com/jwart212/protoc-gen-go-mapper/pkg/types"
)

func TestNullableConverterMatch(t *testing.T) {
	c := NullableConverter{}

	tests := []struct {
		name string
		src  types.TypeInfo
		dst  types.TypeInfo
		want bool
	}{
		{
			name: "Nullable to scalar",
			src:  types.TypeInfo{Kind: types.KindNullable},
			dst:  types.TypeInfo{Kind: types.KindScalar},
			want: true,
		},
		{
			name: "scalar to Nullable",
			src:  types.TypeInfo{Kind: types.KindScalar},
			dst:  types.TypeInfo{Kind: types.KindNullable},
			want: true,
		},
		{
			name: "Nullable to Nullable",
			src:  types.TypeInfo{Kind: types.KindNullable},
			dst:  types.TypeInfo{Kind: types.KindNullable},
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
				t.Errorf("NullableConverter.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullableConverterPriority(t *testing.T) {
	c := NullableConverter{}
	if c.Priority() != 10 {
		t.Errorf("Expected priority 10, got %d", c.Priority())
	}
}

func TestNullableConverterGenerate(t *testing.T) {
	c := NullableConverter{}

	tests := []struct {
		name  string
		field converter.MappingField
		want  string
	}{
		{
			name: "Nullable to scalar",
			field: converter.MappingField{
				SourceExpr: "src",
				TargetExpr: "dst",
				SourceType: types.TypeInfo{Kind: types.KindNullable, Name: "sql.NullString"},
				TargetType: types.TypeInfo{Kind: types.KindScalar, Name: "string"},
			},
			want: "src.String",
		},
		{
			name: "scalar to Nullable",
			field: converter.MappingField{
				SourceExpr: "src",
				TargetExpr: "dst",
				SourceType: types.TypeInfo{Kind: types.KindScalar, Name: "string"},
				TargetType: types.TypeInfo{Kind: types.KindNullable, Name: "sql.NullString"},
			},
			want: "sql.NullString{String: src, Valid: true}",
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

package registry

import (
	"testing"

	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/converter"
	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/types"
)

func TestUUIDConverterMatch(t *testing.T) {
	c := UUIDConverter{}

	tests := []struct {
		name string
		src  types.TypeInfo
		dst  types.TypeInfo
		want bool
	}{
		{
			name: "pgtype.UUID to string",
			src:  types.TypeInfo{Kind: types.KindUUID, Name: "pgtype.UUID"},
			dst:  types.TypeInfo{Kind: types.KindScalar, Name: "string"},
			want: true,
		},
		{
			name: "string to pgtype.UUID",
			src:  types.TypeInfo{Kind: types.KindScalar, Name: "string"},
			dst:  types.TypeInfo{Kind: types.KindUUID, Name: "pgtype.UUID"},
			want: true,
		},
		{
			name: "pgtype.UUID to pgtype.UUID",
			src:  types.TypeInfo{Kind: types.KindUUID, Name: "pgtype.UUID"},
			dst:  types.TypeInfo{Kind: types.KindUUID, Name: "pgtype.UUID"},
			want: true,
		},
		{
			name: "uuid.UUID to string",
			src:  types.TypeInfo{Kind: types.KindUUID, Name: "uuid.UUID"},
			dst:  types.TypeInfo{Kind: types.KindScalar, Name: "string"},
			want: true,
		},
		{
			name: "string to uuid.UUID",
			src:  types.TypeInfo{Kind: types.KindScalar, Name: "string"},
			dst:  types.TypeInfo{Kind: types.KindUUID, Name: "uuid.UUID"},
			want: true,
		},
		{
			name: "scalar to scalar",
			src:  types.TypeInfo{Kind: types.KindScalar, Name: "string"},
			dst:  types.TypeInfo{Kind: types.KindScalar, Name: "string"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Match(tt.src, tt.dst); got != tt.want {
				t.Errorf("UUIDConverter.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUUIDConverterPriority(t *testing.T) {
	c := UUIDConverter{}
	if c.Priority() != 20 {
		t.Errorf("Expected priority 20, got %d", c.Priority())
	}
}

func TestUUIDConverterGenerate(t *testing.T) {
	c := UUIDConverter{}

	tests := []struct {
		name  string
		field converter.MappingField
		want  string
	}{
		{
			name: "uuid.UUID to string",
			field: converter.MappingField{
				SourceExpr: "src",
				TargetExpr: "dst",
				SourceType: types.TypeInfo{Kind: types.KindUUID, Name: "uuid.UUID"},
				TargetType: types.TypeInfo{Kind: types.KindScalar, Name: "string"},
			},
			want: "src.String()",
		},
		{
			name: "string to uuid.UUID",
			field: converter.MappingField{
				SourceExpr: "src",
				TargetExpr: "dst",
				SourceType: types.TypeInfo{Kind: types.KindScalar, Name: "string"},
				TargetType: types.TypeInfo{Kind: types.KindUUID, Name: "uuid.UUID"},
			},
			want: "uuid.MustParse(src)",
		},
		{
			name: "pgtype.UUID to pgtype.UUID",
			field: converter.MappingField{
				SourceExpr: "src",
				TargetExpr: "dst",
				SourceType: types.TypeInfo{Kind: types.KindUUID, Name: "pgtype.UUID"},
				TargetType: types.TypeInfo{Kind: types.KindUUID, Name: "pgtype.UUID"},
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

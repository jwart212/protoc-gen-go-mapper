package registry

import (
	"testing"

	"github.com/jwart212/protoc-gen-go-mapper/pkg/converter"
	"github.com/jwart212/protoc-gen-go-mapper/pkg/types"
)

func TestSliceConverterMatch(t *testing.T) {
	c := SliceConverter{}

	tests := []struct {
		name string
		src  types.TypeInfo
		dst  types.TypeInfo
		want bool
	}{
		{
			name: "slice to slice",
			src:  types.TypeInfo{IsSlice: true},
			dst:  types.TypeInfo{IsSlice: true},
			want: true,
		},
		{
			name: "slice to scalar",
			src:  types.TypeInfo{IsSlice: true},
			dst:  types.TypeInfo{IsSlice: false},
			want: false,
		},
		{
			name: "scalar to slice",
			src:  types.TypeInfo{IsSlice: false},
			dst:  types.TypeInfo{IsSlice: true},
			want: false,
		},
		{
			name: "scalar to scalar",
			src:  types.TypeInfo{IsSlice: false},
			dst:  types.TypeInfo{IsSlice: false},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Match(tt.src, tt.dst); got != tt.want {
				t.Errorf("SliceConverter.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceConverterPriority(t *testing.T) {
	c := SliceConverter{}
	if c.Priority() != 10 {
		t.Errorf("Expected priority 10, got %d", c.Priority())
	}
}

func TestSliceConverterGenerate(t *testing.T) {
	c := SliceConverter{}

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

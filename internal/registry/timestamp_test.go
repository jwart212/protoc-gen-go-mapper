package registry

import (
	"testing"

	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/converter"
	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/types"
)

func TestTimestampConverterMatch(t *testing.T) {
	c := TimestampConverter{}

	tests := []struct {
		name string
		src  types.TypeInfo
		dst  types.TypeInfo
		want bool
	}{
		{
			name: "time.Time to proto Timestamp",
			src:  types.TypeInfo{Kind: types.KindTimestamp, Name: "time.Time"},
			dst:  types.TypeInfo{Kind: types.KindTimestamp, Name: "Timestamp", Package: "timestamppb"},
			want: true,
		},
		{
			name: "proto Timestamp to time.Time",
			src:  types.TypeInfo{Kind: types.KindTimestamp, Name: "Timestamp", Package: "timestamppb"},
			dst:  types.TypeInfo{Kind: types.KindTimestamp, Name: "time.Time"},
			want: true,
		},
		{
			name: "pgtype.Timestamptz to proto Timestamp",
			src:  types.TypeInfo{Kind: types.KindTimestamp, Name: "pgtype.Timestamptz"},
			dst:  types.TypeInfo{Kind: types.KindTimestamp, Name: "Timestamp", Package: "timestamppb"},
			want: true,
		},
		{
			name: "proto Timestamp to pgtype.Timestamptz",
			src:  types.TypeInfo{Kind: types.KindTimestamp, Name: "Timestamp", Package: "timestamppb"},
			dst:  types.TypeInfo{Kind: types.KindTimestamp, Name: "pgtype.Timestamptz"},
			want: true,
		},
		{
			name: "Timestamp to Timestamp (same proto type)",
			src:  types.TypeInfo{Kind: types.KindTimestamp, Name: "Timestamp", Package: "timestamppb"},
			dst:  types.TypeInfo{Kind: types.KindTimestamp, Name: "Timestamp", Package: "timestamppb"},
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
				t.Errorf("TimestampConverter.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimestampConverterPriority(t *testing.T) {
	c := TimestampConverter{}
	if c.Priority() != 20 {
		t.Errorf("Expected priority 20, got %d", c.Priority())
	}
}

func TestTimestampConverterGenerate(t *testing.T) {
	c := TimestampConverter{}

	tests := []struct {
		name  string
		field converter.MappingField
		want  string
	}{
		{
			name: "time.Time to proto Timestamp",
			field: converter.MappingField{
				SourceExpr: "src",
				TargetExpr: "dst",
				SourceType: types.TypeInfo{Kind: types.KindTimestamp, Name: "time.Time"},
				TargetType: types.TypeInfo{Kind: types.KindTimestamp, Name: "Timestamp", Package: "timestamppb"},
			},
			want: "timestamppb.New(src)",
		},
		{
			name: "proto Timestamp to time.Time",
			field: converter.MappingField{
				SourceExpr: "src",
				TargetExpr: "dst",
				SourceType: types.TypeInfo{Kind: types.KindTimestamp, Name: "Timestamp", Package: "timestamppb"},
				TargetType: types.TypeInfo{Kind: types.KindTimestamp, Name: "time.Time"},
			},
			want: "src.AsTime()",
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

package resolver

import (
	"testing"

	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/types"
)

func TestNew(t *testing.T) {
	r := New("sqlc")
	if r == nil {
		t.Error("New() returned nil")
	}
	if r.database != "sqlc" {
		t.Errorf("Expected database to be sqlc, got %s", r.database)
	}
}

func TestResolve(t *testing.T) {
	r := New("sqlc")

	protoType := types.TypeInfo{
		Kind: types.KindUUID,
		Name: "uuid.UUID",
	}

	dbType := r.Resolve(protoType)
	if dbType.Name != "pgtype.UUID" {
		t.Errorf("Expected Name to be pgtype.UUID, got %s", dbType.Name)
	}
}

func TestResolveSQLC(t *testing.T) {
	r := New("sqlc")

	tests := []struct {
		name      string
		protoType types.TypeInfo
		wantName  string
	}{
		{
			name:      "UUID",
			protoType: types.TypeInfo{Kind: types.KindUUID},
			wantName:  "pgtype.UUID",
		},
		{
			name:      "Timestamp",
			protoType: types.TypeInfo{Kind: types.KindTimestamp},
			wantName:  "pgtype.Timestamptz",
		},
		{
			name:      "Scalar",
			protoType: types.TypeInfo{Kind: types.KindScalar, Name: "string"},
			wantName:  "string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.resolveSQLC(tt.protoType, false)
			if got.Name != tt.wantName {
				t.Errorf("resolveSQLC() Name = %v, want %v", got.Name, tt.wantName)
			}
		})
	}
}

func TestResolvePGX(t *testing.T) {
	r := New("pgx")

	tests := []struct {
		name      string
		protoType types.TypeInfo
		wantName  string
	}{
		{
			name:      "UUID",
			protoType: types.TypeInfo{Kind: types.KindUUID},
			wantName:  "pgtype.UUID",
		},
		{
			name:      "Timestamp",
			protoType: types.TypeInfo{Kind: types.KindTimestamp},
			wantName:  "pgtype.Timestamp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.resolvePGX(tt.protoType, false)
			if got.Name != tt.wantName {
				t.Errorf("resolvePGX() Name = %v, want %v", got.Name, tt.wantName)
			}
		})
	}
}

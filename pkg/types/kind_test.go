package types

import "testing"

func TestKindString(t *testing.T) {
	tests := []struct {
		name string
		k    Kind
		want string
	}{
		{"Scalar", KindScalar, "Scalar"},
		{"UUID", KindUUID, "UUID"},
		{"Timestamp", KindTimestamp, "Timestamp"},
		{"Decimal", KindDecimal, "Decimal"},
		{"Enum", KindEnum, "Enum"},
		{"Nullable", KindNullable, "Nullable"},
		{"Message", KindMessage, "Message"},
		{"Unknown", Kind(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.k.String(); got != tt.want {
				t.Errorf("Kind.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

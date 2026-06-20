package generator

import (
	"testing"
)

func TestNewHelper(t *testing.T) {
	h := NewHelper()
	if h == nil {
		t.Error("NewHelper() returned nil")
	}
}

func TestToCamelCase(t *testing.T) {
	h := NewHelper()

	tests := []struct {
		input string
		want  string
	}{
		{"user_name", "userName"},
		{"name", "name"},
		{"", ""},
		{"_private", "Private"},
		{"name_", "name"},
		{"user__name", "userName"},
		{"USER_NAME", "USERName"},
	}

	for _, tt := range tests {
		got := h.ToCamelCase(tt.input)
		if got != tt.want {
			t.Errorf("ToCamelCase(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestToPascalCase(t *testing.T) {
	h := NewHelper()

	tests := []struct {
		input string
		want  string
	}{
		{"user_name", "UserName"},
		{"name", "Name"},
		{"", ""},
		{"_private", "Private"},
		{"name_", "Name"},
		{"user__name", "UserName"},
	}

	for _, tt := range tests {
		got := h.ToPascalCase(tt.input)
		if got != tt.want {
			t.Errorf("ToPascalCase(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestMapSlice(t *testing.T) {
	h := NewHelper()

	got := h.MapSlice("src", "string")
	want := "MapSlice(src, func(item string) string { return item })"
	if got != want {
		t.Errorf("MapSlice() = %q, want %q", got, want)
	}
}

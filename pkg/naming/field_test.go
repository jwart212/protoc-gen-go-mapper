package naming

import "testing"

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"snake_case", "user_name", "userName"},
		{"single word", "name", "name"},
		{"empty", "", ""},
		{"leading underscore", "_private", "Private"},
		{"trailing underscore", "name_", "name"},
		{"multiple underscores", "user__name", "userName"},
		{"SCREAMING_SNAKE", "USER_NAME", "userName"},
		{"already camel", "userName", "username"},
		{"single char", "a", "a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToCamelCase(tt.input); got != tt.want {
				t.Errorf("ToCamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"snake_case", "user_name", "UserName"},
		{"single word", "name", "Name"},
		{"empty", "", ""},
		{"leading underscore", "_private", "Private"},
		{"trailing underscore", "name_", "Name"},
		{"multiple underscores", "user__name", "UserName"},
		{"SCREAMING_SNAKE", "USER_NAME", "UserName"},
		{"already pascal", "UserName", "Username"},
		{"single char", "a", "A"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToPascalCase(tt.input); got != tt.want {
				t.Errorf("ToPascalCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

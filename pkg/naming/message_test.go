package naming

import "testing"

func TestToProtoMessageName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"snake_case", "users", "Users"},
		{"single word", "user", "User"},
		{"empty", "", ""},
		{"compound", "user_profiles", "UserProfiles"},
		{"already pascal", "UserProfile", "Userprofile"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToProtoMessageName(tt.input); got != tt.want {
				t.Errorf("ToProtoMessageName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToDBTableName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"pascal case", "UserProfile", "user_profile"},
		{"single word", "User", "user"},
		{"empty", "", ""},
		{"compound", "UserProfile", "user_profile"},
		{"already snake", "user_profile", "user_profile"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToDBTableName(tt.input); got != tt.want {
				t.Errorf("ToDBTableName() = %v, want %v", got, tt.want)
			}
		})
	}
}

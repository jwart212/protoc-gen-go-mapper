package schema

import "testing"

func TestEnum(t *testing.T) {
	enum := &Enum{
		Name:   "Status",
		Values: []string{"ACTIVE", "INACTIVE", "PENDING"},
	}

	if enum.Name != "Status" {
		t.Errorf("Expected Name to be Status, got %s", enum.Name)
	}
	if len(enum.Values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(enum.Values))
	}
}

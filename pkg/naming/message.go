package naming

// ToProtoMessageName converts a database table name to a protobuf message name.
// It converts snake_case to PascalCase and applies common transformations.
func ToProtoMessageName(table string) string {
	if table == "" {
		return table
	}

	// Convert to PascalCase
	return ToPascalCase(table)
}

// ToDBTableName converts a protobuf message name to a database table name.
// It converts PascalCase to snake_case.
func ToDBTableName(message string) string {
	if message == "" {
		return message
	}

	result := make([]rune, 0, len(message)*2)

	for i, r := range message {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, toLower(r))
	}

	return string(result)
}

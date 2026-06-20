package naming

// ToCamelCase converts a snake_case or SCREAMING_SNAKE_CASE field name to camelCase.
func ToCamelCase(s string) string {
	if s == "" {
		return s
	}

	result := make([]rune, 0, len(s))
	capNext := false

	for i, r := range s {
		if r == '_' {
			capNext = true
			continue
		}

		if i == 0 && capNext {
			// Handle leading underscore
			result = append(result, r)
			capNext = false
			continue
		}

		if capNext {
			result = append(result, toUpper(r))
			capNext = false
		} else {
			result = append(result, toLower(r))
		}
	}

	return string(result)
}

// ToPascalCase converts a snake_case or SCREAMING_SNAKE_CASE field name to PascalCase.
func ToPascalCase(s string) string {
	if s == "" {
		return s
	}

	result := make([]rune, 0, len(s))
	capNext := true

	for _, r := range s {
		if r == '_' {
			capNext = true
			continue
		}

		if capNext {
			result = append(result, toUpper(r))
			capNext = false
		} else {
			result = append(result, toLower(r))
		}
	}

	return string(result)
}

func toUpper(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - ('a' - 'A')
	}
	return r
}

func toLower(r rune) rune {
	if r >= 'A' && r <= 'Z' {
		return r + ('a' - 'A')
	}
	return r
}

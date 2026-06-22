package converter

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ConvertUUID converts pgtype.UUID to either string or *string based on type parameter T.
// This is a generic converter that handles both nullable and non-nullable UUID fields.
func ConvertUUID[T string | *string](v pgtype.UUID) T {
	if v.Valid {
		s := uuid.UUID(v.Bytes).String()
		// Type assertion to determine if T is string or *string
		var t T
		switch any(t).(type) {
		case string:
			return any(s).(T)
		case *string:
			return any(&s).(T)
		}
	}
	var zero T
	return zero
}

// ConvertTimestamp converts pgtype.Timestamptz to *timestamppb.Timestamp.
// This is a generic converter for timestamp fields.
func ConvertTimestamp[T *timestamppb.Timestamp](v pgtype.Timestamptz) T {
	if v.Valid {
		return timestamppb.New(v.Time)
	}
	var zero T
	return zero
}

// ConvertText converts pgtype.Text to either string or *string based on type parameter T.
// This is a generic converter that handles both nullable and non-nullable text fields.
func ConvertText[T string | *string](v pgtype.Text) T {
	if v.Valid {
		// Type assertion to determine if T is string or *string
		var t T
		switch any(t).(type) {
		case string:
			return any(v.String).(T)
		case *string:
			return any(&v.String).(T)
		}
	}
	var zero T
	return zero
}

// ConvertInt32 converts int32 to *int32 for nullable fields.
func ConvertInt32[T int32 | *int32](v int32, valid bool) T {
	if valid {
		var t T
		switch any(t).(type) {
		case int32:
			return any(v).(T)
		case *int32:
			return any(&v).(T)
		}
	}
	var zero T
	return zero
}

// ConvertInt64 converts int64 to *int64 for nullable fields.
func ConvertInt64[T int64 | *int64](v int64, valid bool) T {
	if valid {
		var t T
		switch any(t).(type) {
		case int64:
			return any(v).(T)
		case *int64:
			return any(&v).(T)
		}
	}
	var zero T
	return zero
}

// ConvertBool converts bool to *bool for nullable fields.
func ConvertBool[T bool | *bool](v bool, valid bool) T {
	if valid {
		var t T
		switch any(t).(type) {
		case bool:
			return any(v).(T)
		case *bool:
			return any(&v).(T)
		}
	}
	var zero T
	return zero
}

// ConvertFloat64 converts float64 to *float64 for nullable fields.
func ConvertFloat64[T float64 | *float64](v float64, valid bool) T {
	if valid {
		var t T
		switch any(t).(type) {
		case float64:
			return any(v).(T)
		case *float64:
			return any(&v).(T)
		}
	}
	var zero T
	return zero
}

// ConvertStringToNumeric converts string to numeric types (int32, int64, float64).
// This is useful for converting string fields from proto to numeric DB fields.
// Returns zero value if parsing fails.
func ConvertStringToNumeric[T int32 | int64 | float64](v string) T {
	var result T
	switch any(result).(type) {
	case int32:
		var n int32
		if _, err := fmt.Sscanf(v, "%d", &n); err == nil {
			result = any(n).(T)
		}
	case int64:
		var n int64
		if _, err := fmt.Sscanf(v, "%d", &n); err == nil {
			result = any(n).(T)
		}
	case float64:
		var n float64
		if _, err := fmt.Sscanf(v, "%f", &n); err == nil {
			result = any(n).(T)
		}
	}
	return result
}

// ConvertStringToNumericPtr converts string to nullable numeric types (*int32, *int64, *float64).
// Returns nil if parsing fails.
func ConvertStringToNumericPtr[T *int32 | *int64 | *float64](v string) T {
	switch any(T(nil)).(type) {
	case *int32:
		var n int32
		if _, err := fmt.Sscanf(v, "%d", &n); err == nil {
			return any(&n).(T)
		}
		return any((*int32)(nil)).(T)
	case *int64:
		var n int64
		if _, err := fmt.Sscanf(v, "%d", &n); err == nil {
			return any(&n).(T)
		}
		return any((*int64)(nil)).(T)
	case *float64:
		var n float64
		if _, err := fmt.Sscanf(v, "%f", &n); err == nil {
			return any(&n).(T)
		}
		return any((*float64)(nil)).(T)
	}
	var zero T
	return zero
}

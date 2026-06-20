package converter

import (
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
	"github.com/google/uuid"
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

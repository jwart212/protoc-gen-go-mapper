package registry

import (
	"fmt"

	"github.com/jwart212/protoc-gen-go-mapper/pkg/converter"
	"github.com/jwart212/protoc-gen-go-mapper/pkg/types"
)

// NullableConverter handles sql.Null* types ↔ optional field conversions.
type NullableConverter struct{}

// Match returns true for Nullable ↔ scalar conversions (excluding timestamps and UUIDs).
func (c NullableConverter) Match(src, dst types.TypeInfo) bool {
	// Don't match if target is a timestamp type (let TimestampConverter handle it)
	if dst.Kind == types.KindTimestamp && dst.Name == "google.protobuf.Timestamp" {
		return false
	}
	if src.Kind == types.KindTimestamp && src.Name == "google.protobuf.Timestamp" {
		return false
	}
	// Don't match if either side is UUID type (let UUIDConverter handle it)
	if dst.Kind == types.KindUUID || dst.Name == "pgtype.UUID" {
		return false
	}
	if src.Kind == types.KindUUID || src.Name == "pgtype.UUID" {
		return false
	}
	// Nullable to scalar
	if src.Kind == types.KindNullable && dst.Kind == types.KindScalar {
		return true
	}
	// scalar to Nullable
	if src.Kind == types.KindScalar && dst.Kind == types.KindNullable {
		return true
	}
	return false
}

// Priority returns a higher priority than ScalarConverter for Nullable-specific conversions.
func (c NullableConverter) Priority() int {
	return 10
}

// Generate emits the Go expression for Nullable conversion.
func (c NullableConverter) Generate(field converter.MappingField) (string, error) {
	// Nullable to pointer (DB → Proto)
	if field.SourceType.Kind == types.KindNullable && field.TargetType.Kind == types.KindScalar && field.TargetType.IsNullable {
		switch field.SourceType.Name {
		case "sql.NullInt32":
			return fmt.Sprintf("newInt32(%s.Int32, %s.Valid)", field.SourceExpr, field.SourceExpr), nil
		case "sql.NullInt64":
			return fmt.Sprintf("newInt64(%s.Int64, %s.Valid)", field.SourceExpr, field.SourceExpr), nil
		case "sql.NullBool":
			return fmt.Sprintf("newBool(%s.Bool, %s.Valid)", field.SourceExpr, field.SourceExpr), nil
		case "sql.NullString":
			return fmt.Sprintf("newString(%s.String, %s.Valid)", field.SourceExpr, field.SourceExpr), nil
		case "sql.NullFloat64":
			return fmt.Sprintf("newFloat64(%s.Float64, %s.Valid)", field.SourceExpr, field.SourceExpr), nil
		case "sql.NullTime":
			// For time, convert to string in proto
			return fmt.Sprintf("newTimeString(%s.Time, %s.Valid)", field.SourceExpr, field.SourceExpr), nil
		default:
			return fmt.Sprintf("newString(%s.String, %s.Valid)", field.SourceExpr, field.SourceExpr), nil
		}
	}
	// Nullable to scalar/timestamp (DB → Proto)
	if field.SourceType.Kind == types.KindNullable && (field.TargetType.Kind == types.KindScalar || field.TargetType.Kind == types.KindTimestamp) {
		switch field.SourceType.Name {
		case "sql.NullInt32":
			return fmt.Sprintf("%s.Int32", field.SourceExpr), nil
		case "sql.NullInt64":
			return fmt.Sprintf("%s.Int64", field.SourceExpr), nil
		case "sql.NullBool":
			return fmt.Sprintf("%s.Bool", field.SourceExpr), nil
		case "sql.NullString":
			return fmt.Sprintf("%s.String", field.SourceExpr), nil
		case "sql.NullFloat64":
			return fmt.Sprintf("%s.Float64", field.SourceExpr), nil
		case "sql.NullTime":
			// For time, convert to string - check if target is nullable (pointer)
			if field.TargetType.IsNullable {
				return fmt.Sprintf("newTimeString(%s.Time, %s.Valid)", field.SourceExpr, field.SourceExpr), nil
			}
			return fmt.Sprintf("%s.Time.Format(\"2006-01-02\")", field.SourceExpr), nil
		default:
			return fmt.Sprintf("%s.String", field.SourceExpr), nil
		}
	}
	// pointer to Nullable (Proto → DB)
	if field.SourceType.Kind == types.KindScalar && field.SourceType.IsNullable && field.TargetType.Kind == types.KindNullable {
		switch field.TargetType.Name {
		case "sql.NullInt32":
			return fmt.Sprintf("newNullInt32(%s)", field.SourceExpr), nil
		case "sql.NullInt64":
			return fmt.Sprintf("newNullInt64(%s)", field.SourceExpr), nil
		case "sql.NullBool":
			return fmt.Sprintf("newNullBool(%s)", field.SourceExpr), nil
		case "sql.NullString":
			return fmt.Sprintf("newNullString(%s)", field.SourceExpr), nil
		case "sql.NullFloat64":
			return fmt.Sprintf("newNullFloat64(%s)", field.SourceExpr), nil
		case "sql.NullTime":
			// For time, parse from string
			return fmt.Sprintf("newNullTimeFromString(%s)", field.SourceExpr), nil
		default:
			return fmt.Sprintf("newNullString(%s)", field.SourceExpr), nil
		}
	}
	// timestamp to Nullable (Proto → DB)
	if field.SourceType.Kind == types.KindTimestamp && field.SourceType.IsNullable && field.TargetType.Kind == types.KindNullable {
		switch field.TargetType.Name {
		case "sql.NullTime":
			return fmt.Sprintf("newNullTimeFromString(%s)", field.SourceExpr), nil
		default:
			return fmt.Sprintf("newNullString(%s)", field.SourceExpr), nil
		}
	}
	// scalar to Nullable (Proto → DB)
	if field.SourceType.Kind == types.KindScalar && !field.SourceType.IsNullable && field.TargetType.Kind == types.KindNullable {
		switch field.TargetType.Name {
		case "sql.NullInt32":
			return fmt.Sprintf("sql.NullInt32{Int32: %s, Valid: true}", field.SourceExpr), nil
		case "sql.NullInt64":
			return fmt.Sprintf("sql.NullInt64{Int64: %s, Valid: true}", field.SourceExpr), nil
		case "sql.NullBool":
			return fmt.Sprintf("sql.NullBool{Bool: %s, Valid: true}", field.SourceExpr), nil
		case "sql.NullString":
			return fmt.Sprintf("sql.NullString{String: %s, Valid: true}", field.SourceExpr), nil
		case "sql.NullFloat64":
			return fmt.Sprintf("sql.NullFloat64{Float64: %s, Valid: true}", field.SourceExpr), nil
		case "sql.NullTime":
			return fmt.Sprintf("sql.NullTime{Time: %s, Valid: true}", field.SourceExpr), nil
		default:
			return fmt.Sprintf("sql.NullString{String: %s, Valid: true}", field.SourceExpr), nil
		}
	}
	return "", fmt.Errorf("unsupported nullable conversion: %v -> %v", field.SourceType, field.TargetType)
}

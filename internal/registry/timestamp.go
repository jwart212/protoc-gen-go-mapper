package registry

import (
	"fmt"

	"github.com/jwart212/protoc-gen-go-mapper/pkg/converter"
	"github.com/jwart212/protoc-gen-go-mapper/pkg/types"
)

// TimestampConverter handles time.Time ↔ google.protobuf.Timestamp conversions.
type TimestampConverter struct{}

// Match returns true for Timestamp ↔ time.Time and Timestamp ↔ nullable conversions.
func (c TimestampConverter) Match(src, dst types.TypeInfo) bool {
	// Helper to check if a type is google.protobuf.Timestamp
	isProtoTimestamp := func(t types.TypeInfo) bool {
		return t.Kind == types.KindTimestamp && t.Name == "Timestamp" && t.Package == "timestamppb"
	}

	// time.Time to google.protobuf.Timestamp
	if src.Kind == types.KindTimestamp && src.Name == "time.Time" && isProtoTimestamp(dst) {
		return true
	}
	// google.protobuf.Timestamp to time.Time
	if isProtoTimestamp(src) && dst.Kind == types.KindTimestamp && dst.Name == "time.Time" {
		return true
	}
	// pgtype.Timestamptz to google.protobuf.Timestamp
	if src.Kind == types.KindTimestamp && src.Name == "pgtype.Timestamptz" && isProtoTimestamp(dst) {
		return true
	}
	// google.protobuf.Timestamp to pgtype.Timestamptz
	if isProtoTimestamp(src) && dst.Kind == types.KindTimestamp && dst.Name == "pgtype.Timestamptz" {
		return true
	}
	// google.protobuf.Timestamp (nullable) to pgtype.Timestamptz (nullable)
	if isProtoTimestamp(src) && src.IsNullable && dst.Kind == types.KindNullable && dst.Name == "pgtype.Timestamptz" {
		return true
	}
	// pgtype.Timestamptz (nullable) to google.protobuf.Timestamp (nullable)
	if src.Kind == types.KindNullable && src.Name == "pgtype.Timestamptz" && isProtoTimestamp(dst) && dst.IsNullable {
		return true
	}
	// sql.NullTime to google.protobuf.Timestamp
	if src.Kind == types.KindNullable && src.Name == "sql.NullTime" && isProtoTimestamp(dst) {
		return true
	}
	// google.protobuf.Timestamp to sql.NullTime
	if isProtoTimestamp(src) && dst.Kind == types.KindNullable && dst.Name == "sql.NullTime" {
		return true
	}
	return false
}

// Priority returns a higher priority than NullableConverter for Timestamp-specific conversions.
func (c TimestampConverter) Priority() int {
	return 20
}

// Generate emits the Go expression for Timestamp conversion.
func (c TimestampConverter) Generate(field converter.MappingField) (string, error) {
	// Helper to check if a type is google.protobuf.Timestamp
	isProtoTimestamp := func(t types.TypeInfo) bool {
		return t.Kind == types.KindTimestamp && t.Name == "Timestamp" && t.Package == "timestamppb"
	}

	// time.Time to google.protobuf.Timestamp
	if field.SourceType.Kind == types.KindTimestamp && field.SourceType.Name == "time.Time" && isProtoTimestamp(field.TargetType) {
		return fmt.Sprintf("timestamppb.New(%s)", field.SourceExpr), nil
	}
	// google.protobuf.Timestamp to time.Time
	if isProtoTimestamp(field.SourceType) && field.TargetType.Kind == types.KindTimestamp && field.TargetType.Name == "time.Time" {
		return fmt.Sprintf("%s.AsTime()", field.SourceExpr), nil
	}
	// pgtype.Timestamptz to google.protobuf.Timestamp
	if field.SourceType.Kind == types.KindTimestamp && field.SourceType.Name == "pgtype.Timestamptz" && isProtoTimestamp(field.TargetType) {
		return fmt.Sprintf("newTimestampFromTimestamptz(%s)", field.SourceExpr), nil
	}
	// google.protobuf.Timestamp to pgtype.Timestamptz
	if isProtoTimestamp(field.SourceType) && field.TargetType.Kind == types.KindTimestamp && field.TargetType.Name == "pgtype.Timestamptz" {
		return fmt.Sprintf("newTimestamptzFromTimestamp(%s)", field.SourceExpr), nil
	}
	// google.protobuf.Timestamp (nullable) to pgtype.Timestamptz (nullable)
	if isProtoTimestamp(field.SourceType) && field.SourceType.IsNullable && field.TargetType.Kind == types.KindNullable && field.TargetType.Name == "pgtype.Timestamptz" {
		return fmt.Sprintf("newTimestamptzFromTimestamp(%s)", field.SourceExpr), nil
	}
	// pgtype.Timestamptz (nullable) to google.protobuf.Timestamp (nullable)
	if field.SourceType.Kind == types.KindNullable && field.SourceType.Name == "pgtype.Timestamptz" && isProtoTimestamp(field.TargetType) && field.TargetType.IsNullable {
		return fmt.Sprintf("newTimestampFromTimestamptz(%s)", field.SourceExpr), nil
	}
	// sql.NullTime to google.protobuf.Timestamp (nullable)
	if field.SourceType.Kind == types.KindNullable && field.SourceType.Name == "sql.NullTime" && isProtoTimestamp(field.TargetType) {
		return fmt.Sprintf("newTimestampFromNullTime(%s)", field.SourceExpr), nil
	}
	// google.protobuf.Timestamp to sql.NullTime (nullable)
	if isProtoTimestamp(field.SourceType) && field.TargetType.Kind == types.KindNullable && field.TargetType.Name == "sql.NullTime" {
		return fmt.Sprintf("newNullTimeFromTimestamp(%s)", field.SourceExpr), nil
	}
	return "", fmt.Errorf("unsupported timestamp conversion: %v -> %v", field.SourceType, field.TargetType)
}

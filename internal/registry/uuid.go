package registry

import (
	"fmt"

	"github.com/jwart212/protoc-gen-go-mapper/pkg/converter"
	"github.com/jwart212/protoc-gen-go-mapper/pkg/types"
)

// UUIDConverter handles uuid.UUID ↔ string conversions.
type UUIDConverter struct{}

// Match returns true for UUID ↔ scalar (string) and UUID ↔ UUID conversions.
func (c UUIDConverter) Match(src, dst types.TypeInfo) bool {
	// pgtype.UUID to string (check this first before generic UUID to string)
	if src.Kind == types.KindUUID && src.Name == "pgtype.UUID" && dst.Kind == types.KindScalar && dst.Name == "string" {
		return true
	}
	// string to pgtype.UUID (check this first before generic string to UUID)
	if src.Kind == types.KindScalar && src.Name == "string" && dst.Kind == types.KindUUID && dst.Name == "pgtype.UUID" {
		return true
	}
	// pgtype.UUID to pgtype.UUID
	if src.Kind == types.KindUUID && src.Name == "pgtype.UUID" && dst.Kind == types.KindUUID && dst.Name == "pgtype.UUID" {
		return true
	}
	// nullable pgtype.UUID to string (for optional ID fields)
	if src.Kind == types.KindNullable && src.Name == "pgtype.UUID" && dst.Kind == types.KindScalar && dst.Name == "string" {
		return true
	}
	// string to nullable pgtype.UUID (for optional ID fields)
	if src.Kind == types.KindScalar && src.Name == "string" && dst.Kind == types.KindNullable && dst.Name == "pgtype.UUID" {
		return true
	}
	// pgtype.UUID to nullable string (for DB UUID to optional proto string)
	if src.Kind == types.KindUUID && src.Name == "pgtype.UUID" && dst.Kind == types.KindScalar && dst.Name == "string" && dst.IsNullable {
		return true
	}
	// string to pgtype.UUID (optional string to non-nullable UUID)
	if src.Kind == types.KindScalar && src.Name == "string" && src.IsNullable && dst.Kind == types.KindUUID && dst.Name == "pgtype.UUID" && !dst.IsNullable {
		return true
	}
	// UUID to string (generic, for uuid.UUID to string)
	if src.Kind == types.KindUUID && dst.Kind == types.KindScalar && dst.Name == "string" {
		return true
	}
	// string to UUID (generic, for string to uuid.UUID)
	if src.Kind == types.KindScalar && src.Name == "string" && dst.Kind == types.KindUUID {
		return true
	}
	// UUID to UUID (both sides are UUID kind) - handle string to uuid.UUID and uuid.UUID to string
	if src.Kind == types.KindUUID && dst.Kind == types.KindUUID {
		return true
	}
	// UUID to nullable UUID
	if src.Kind == types.KindUUID && dst.Kind == types.KindNullable && dst.Name == "uuid.NullUUID" {
		return true
	}
	// nullable UUID to UUID
	if src.Kind == types.KindNullable && src.Name == "uuid.NullUUID" && dst.Kind == types.KindUUID {
		return true
	}
	// string to nullable UUID
	if src.Kind == types.KindScalar && src.IsNullable && dst.Kind == types.KindNullable && dst.Name == "uuid.NullUUID" {
		return true
	}
	// nullable UUID to string
	if src.Kind == types.KindNullable && src.Name == "uuid.NullUUID" && dst.Kind == types.KindScalar && dst.IsNullable {
		return true
	}
	return false
}

// Priority returns a higher priority than NullableConverter for UUID-specific conversions.
func (c UUIDConverter) Priority() int {
	return 20
}

// Generate emits the Go expression for UUID conversion.
func (c UUIDConverter) Generate(field converter.MappingField) (string, error) {
	// pgtype.UUID to string (non-pointer for non-optional proto fields)
	if field.SourceType.Kind == types.KindUUID && field.SourceType.Name == "pgtype.UUID" && field.TargetType.Kind == types.KindScalar && field.TargetType.Name == "string" && !field.TargetType.IsNullable {
		return fmt.Sprintf("newStringFromUUIDNonPtr(%s)", field.SourceExpr), nil
	}
	// nullable pgtype.UUID to string (pointer for optional proto fields)
	if field.SourceType.Kind == types.KindNullable && field.SourceType.Name == "pgtype.UUID" && field.TargetType.Kind == types.KindScalar && field.TargetType.Name == "string" && field.TargetType.IsNullable {
		return fmt.Sprintf("newStringFromUUID(%s)", field.SourceExpr), nil
	}
	// pgtype.UUID to string (pointer for optional proto fields - when source is KindUUID but nullable)
	if field.SourceType.Kind == types.KindUUID && field.SourceType.Name == "pgtype.UUID" && field.TargetType.Kind == types.KindScalar && field.TargetType.Name == "string" && field.TargetType.IsNullable {
		return fmt.Sprintf("newStringFromUUID(%s)", field.SourceExpr), nil
	}
	// pgtype.UUID to string (non-nullable DB to optional proto field)
	if field.SourceType.Kind == types.KindUUID && field.SourceType.Name == "pgtype.UUID" && !field.SourceType.IsNullable && field.TargetType.Kind == types.KindScalar && field.TargetType.Name == "string" && field.TargetType.IsNullable {
		return fmt.Sprintf("newStringFromUUID(%s)", field.SourceExpr), nil
	}
	// string to pgtype.UUID (non-pointer to non-nullable DB field)
	if field.SourceType.Kind == types.KindScalar && field.SourceType.Name == "string" && !field.SourceType.IsNullable && field.TargetType.Kind == types.KindUUID && field.TargetType.Name == "pgtype.UUID" && !field.TargetType.IsNullable {
		return fmt.Sprintf("newUUIDFromString(%s)", field.SourceExpr), nil
	}
	// string to pgtype.UUID (pointer to nullable DB field)
	if field.SourceType.Kind == types.KindScalar && field.SourceType.Name == "string" && field.SourceType.IsNullable && field.TargetType.Kind == types.KindNullable && field.TargetType.Name == "pgtype.UUID" {
		return fmt.Sprintf("newUUID(%s)", field.SourceExpr), nil
	}
	// string to pgtype.UUID (pointer to non-nullable DB field - for optional proto ID fields to non-nullable DB UUID)
	if field.SourceType.Kind == types.KindScalar && field.SourceType.Name == "string" && field.SourceType.IsNullable && field.TargetType.Kind == types.KindUUID && field.TargetType.Name == "pgtype.UUID" && !field.TargetType.IsNullable {
		return fmt.Sprintf("newUUID(%s)", field.SourceExpr), nil
	}
	// string to pgtype.UUID (pointer to non-nullable DB field - for optional proto string fields like deleted_by)
	if field.SourceType.Kind == types.KindScalar && field.SourceType.Name == "string" && field.SourceType.IsNullable && field.TargetType.Kind == types.KindUUID && field.TargetType.Name == "pgtype.UUID" {
		return fmt.Sprintf("newUUID(%s)", field.SourceExpr), nil
	}
	// UUID to string (generic, for uuid.UUID to string)
	if field.SourceType.Kind == types.KindUUID && field.TargetType.Kind == types.KindScalar && field.TargetType.Name == "string" {
		return fmt.Sprintf("%s.String()", field.SourceExpr), nil
	}
	// string to UUID (generic, for string to uuid.UUID)
	if field.SourceType.Kind == types.KindScalar && field.TargetType.Kind == types.KindUUID {
		return fmt.Sprintf("uuid.MustParse(%s)", field.SourceExpr), nil
	}
	// UUID to UUID (string to uuid.UUID)
	if field.SourceType.Kind == types.KindUUID && field.SourceType.Name == "string" && field.TargetType.Kind == types.KindUUID && field.TargetType.Name == "uuid.UUID" {
		return fmt.Sprintf("uuid.MustParse(%s)", field.SourceExpr), nil
	}
	// UUID to UUID (uuid.UUID to string)
	if field.SourceType.Kind == types.KindUUID && field.SourceType.Name == "uuid.UUID" && field.TargetType.Kind == types.KindUUID && field.TargetType.Name == "string" {
		return fmt.Sprintf("%s.String()", field.SourceExpr), nil
	}
	// pgtype.UUID to pgtype.UUID
	if field.SourceType.Kind == types.KindUUID && field.SourceType.Name == "pgtype.UUID" && field.TargetType.Kind == types.KindUUID && field.TargetType.Name == "pgtype.UUID" {
		return field.SourceExpr, nil
	}
	// UUID to nullable UUID
	if field.SourceType.Kind == types.KindUUID && field.TargetType.Kind == types.KindNullable && field.TargetType.Name == "uuid.NullUUID" {
		return fmt.Sprintf("uuid.NullUUID{UUID: %s, Valid: true}", field.SourceExpr), nil
	}
	// nullable UUID to UUID
	if field.SourceType.Kind == types.KindNullable && field.SourceType.Name == "uuid.NullUUID" && field.TargetType.Kind == types.KindUUID {
		return fmt.Sprintf("%s.UUID", field.SourceExpr), nil
	}
	// string to nullable UUID
	if field.SourceType.Kind == types.KindScalar && field.SourceType.IsNullable && field.TargetType.Kind == types.KindNullable && field.TargetType.Name == "uuid.NullUUID" {
		return fmt.Sprintf("newNullUUID(%s)", field.SourceExpr), nil
	}
	// nullable UUID to string
	if field.SourceType.Kind == types.KindNullable && field.SourceType.Name == "uuid.NullUUID" && field.TargetType.Kind == types.KindScalar && field.TargetType.IsNullable {
		return fmt.Sprintf("newStringFromUUID(%s)", field.SourceExpr), nil
	}
	return "", fmt.Errorf("unsupported UUID conversion: %v -> %v", field.SourceType, field.TargetType)
}

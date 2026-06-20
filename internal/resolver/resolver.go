package resolver

import (
	"strings"

	"github.com/jwart212/protoc-gen-go-mapper/pkg/types"
)

// Resolver maps protobuf types to database-specific types.
type Resolver struct {
	database string
}

// New creates a new Resolver for the specified database type.
func New(database string) *Resolver {
	return &Resolver{
		database: database,
	}
}

// Resolve maps a protobuf TypeInfo to a database TypeInfo.
func (r *Resolver) Resolve(protoType types.TypeInfo) types.TypeInfo {
	dbType := protoType

	// Apply database-specific mappings
	switch r.database {
	case "sqlc":
		dbType = r.resolveSQLC(protoType, false)
	case "pgx":
		dbType = r.resolvePGX(protoType, false)
	case "database_sql":
		dbType = r.resolveDatabaseSQL(protoType)
	}

	return dbType
}

// ResolveWithFieldName maps a protobuf TypeInfo to a database TypeInfo with field name context.
func (r *Resolver) ResolveWithFieldName(protoType types.TypeInfo, fieldName string) types.TypeInfo {
	dbType := protoType

	// Special handling for ID fields - map string to pgtype.UUID for SQLC/PGX
	// Note: We always map to non-nullable pgtype.UUID for ID fields in the DB model
	// The proto side can be optional (pointer) or required (non-pointer)
	isIDField := (fieldName == "id" || strings.HasSuffix(strings.ToLower(fieldName), "_id")) && protoType.Kind == types.KindScalar && protoType.Name == "string"
	if isIDField && (r.database == "sqlc" || r.database == "pgx") {
		// Always map to non-nullable pgtype.UUID for DB side
		dbType.Kind = types.KindUUID
		dbType.Name = "pgtype.UUID"
	}

	// Special handling for deleted_by field - map to pgtype.UUID for SQLC/PGX
	// This is a UUID field that doesn't follow the _id naming convention
	isDeletedByField := fieldName == "deleted_by" && protoType.Kind == types.KindScalar && protoType.Name == "string"
	if isDeletedByField && (r.database == "sqlc" || r.database == "pgx") {
		// Map to non-nullable pgtype.UUID for DB side
		dbType.Kind = types.KindUUID
		dbType.Name = "pgtype.UUID"
	}

	// Apply database-specific mappings
	switch r.database {
	case "sqlc":
		dbType = r.resolveSQLC(dbType, isIDField || isDeletedByField)
	case "pgx":
		dbType = r.resolvePGX(dbType, isIDField || isDeletedByField)
	case "database_sql":
		dbType = r.resolveDatabaseSQL(dbType)
	}

	return dbType
}

// resolveSQLC maps protobuf types to sqlc types.
func (r *Resolver) resolveSQLC(protoType types.TypeInfo, isIDField bool) types.TypeInfo {
	dbType := protoType

	// SQLC-specific type mappings (using pgtype types for PostgreSQL)
	switch protoType.Kind {
	case types.KindUUID:
		// Keep as UUID for DB, but mark as nullable if proto is optional
		// Don't convert to nullable for ID fields even if proto is optional
		if protoType.IsNullable && !isIDField {
			dbType.Kind = types.KindNullable
			dbType.Name = "pgtype.UUID"
		} else {
			dbType.Name = "pgtype.UUID"
		}
	case types.KindTimestamp:
		// For google.protobuf.Timestamp, map to pgtype.Timestamptz
		if protoType.IsNullable {
			dbType.Kind = types.KindNullable
			dbType.Name = "pgtype.Timestamptz"
		} else {
			dbType.Name = "pgtype.Timestamptz"
		}
	case types.KindDecimal:
		dbType.Name = "pgtype.Numeric"
	}

	// Handle nullable types for SQLC
	// Only convert to nullable if the proto type is explicitly marked as optional
	// Skip for ID fields (already handled above)
	if protoType.Kind == types.KindScalar && !protoType.IsSlice && protoType.IsNullable && !isIDField {
		// Map scalar types to nullable types for SQLC
		switch protoType.Name {
		case "int32", "int64", "int":
			dbType.Kind = types.KindNullable
			dbType.Name = "pgtype.Int8"
		case "bool":
			dbType.Kind = types.KindNullable
			dbType.Name = "pgtype.Bool"
		case "string":
			dbType.Kind = types.KindNullable
			dbType.Name = "pgtype.Text"
		case "float64":
			dbType.Kind = types.KindNullable
			dbType.Name = "pgtype.Numeric"
		}
	}

	return dbType
}

// resolvePGX maps protobuf types to pgx types.
func (r *Resolver) resolvePGX(protoType types.TypeInfo, isIDField bool) types.TypeInfo {
	dbType := protoType

	// PGX-specific type mappings
	switch protoType.Kind {
	case types.KindUUID:
		// Keep as UUID for DB, but mark as nullable if proto is optional
		// Don't convert to nullable for ID fields even if proto is optional
		if protoType.IsNullable && !isIDField {
			dbType.Kind = types.KindNullable
			dbType.Name = "pgtype.UUID"
		} else {
			dbType.Name = "pgtype.UUID"
		}
	case types.KindTimestamp:
		dbType.Name = "pgtype.Timestamp"
	case types.KindDecimal:
		dbType.Name = "pgtype.Numeric"
	}

	return dbType
}

// resolveDatabaseSQL maps protobuf types to database/sql types.
func (r *Resolver) resolveDatabaseSQL(protoType types.TypeInfo) types.TypeInfo {
	dbType := protoType

	// database/sql-specific type mappings
	switch protoType.Kind {
	case types.KindUUID:
		dbType.Name = "string"
	case types.KindTimestamp:
		dbType.Name = "time.Time"
	case types.KindDecimal:
		dbType.Name = "string"
	}

	return dbType
}

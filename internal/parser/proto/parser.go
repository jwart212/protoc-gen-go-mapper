package proto

import (
	"fmt"
	"strings"

	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/internal/resolver"
	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/internal/schema"
	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/types"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Parser converts protobuf descriptors into the internal schema model.
type Parser struct {
	resolver *resolver.Resolver
}

// New creates a new Parser instance.
func New(database string) *Parser {
	return &Parser{
		resolver: resolver.New(database),
	}
}

// ParseFile converts a FileDescriptorProto into a schema.Model.
func (p *Parser) ParseFile(fileProto *descriptorpb.FileDescriptorProto) (*schema.Model, error) {
	model := &schema.Model{
		Messages: make([]*schema.Message, 0),
		Enums:    make([]*schema.Enum, 0),
	}

	// Parse messages (skip well-known types and imported messages)
	for _, msgDesc := range fileProto.MessageType {
		// Skip well-known protobuf types
		msgName := msgDesc.GetName()
		if isWellKnownType(msgName) {
			continue
		}
		// Skip messages from other packages by checking if the message name contains a dot
		// Messages defined in the current file have simple names (no dots)
		// Messages from other files have fully qualified names like "google.protobuf.Timestamp"
		if strings.Contains(msgName, ".") {
			continue
		}
		msg, err := p.parseMessage(msgDesc, fileProto)
		if err != nil {
			return nil, fmt.Errorf("parsing message %s: %w", msgDesc.GetName(), err)
		}
		model.Messages = append(model.Messages, msg)
	}

	// Parse enums
	for _, enumDesc := range fileProto.EnumType {
		enum := p.parseEnum(enumDesc)
		model.Enums = append(model.Enums, enum)
	}

	return model, nil
}

// parseMessage converts a DescriptorProto to a schema.Message.
func (p *Parser) parseMessage(msgDesc *descriptorpb.DescriptorProto, fileProto *descriptorpb.FileDescriptorProto) (*schema.Message, error) {
	// Extract package name from go_package option
	// Format is typically "path/to/package;package_name" or just "path/to/package"
	goPackage := fileProto.GetOptions().GetGoPackage()
	var pkgName string
	if goPackage != "" {
		// Split on semicolon to get the package name if specified
		parts := strings.Split(goPackage, ";")
		if len(parts) > 1 {
			pkgName = parts[1]
		} else {
			// Extract the last component of the path as the package name
			pathParts := strings.Split(parts[0], "/")
			pkgName = pathParts[len(pathParts)-1]
		}
	}

	msg := &schema.Message{
		Name:    msgDesc.GetName(),
		Package: pkgName,
		Fields:  make([]*schema.Field, 0),
	}

	for _, fieldDesc := range msgDesc.Field {
		field, err := p.parseField(fieldDesc, fileProto)
		if err != nil {
			return nil, fmt.Errorf("parsing field %s: %w", fieldDesc.GetName(), err)
		}
		msg.Fields = append(msg.Fields, field)
	}

	return msg, nil
}

// parseField converts a FieldDescriptorProto to a schema.Field.
func (p *Parser) parseField(fieldDesc *descriptorpb.FieldDescriptorProto, fileProto *descriptorpb.FileDescriptorProto) (*schema.Field, error) {
	protoType := p.typeInfoFromProtoType(fieldDesc.GetType(), fieldDesc.GetTypeName())

	// Mark as slice if repeated
	repeated := fieldDesc.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED
	if repeated {
		protoType.IsSlice = true
	}

	// Mark as nullable if proto3 optional
	optional := fieldDesc.GetProto3Optional()
	if optional {
		protoType.IsNullable = true
	}

	// Store the original field name for type resolution
	fieldName := fieldDesc.GetName()

	// Special handling for google.protobuf.Timestamp fields
	if fieldDesc.GetTypeName() == ".google.protobuf.Timestamp" {
		protoType.Kind = types.KindTimestamp
		protoType.Name = "Timestamp"
		protoType.Package = "timestamppb"
	}

	// Resolve DB type with field name context (for ID field UUID mapping)
	dbType := p.resolver.ResolveWithFieldName(protoType, fieldName)

	// Clear package on DB type - it should not inherit proto package
	dbType.Package = ""

	// Copy nullable flag to DB type
	if optional {
		dbType.IsNullable = true
	}

	// Copy slice flag to DB type
	if repeated {
		dbType.IsSlice = true
	}

	return &schema.Field{
		Name:        fieldName,
		ProtoType:   protoType,
		DBType:      dbType,
		FieldNumber: fieldDesc.GetNumber(),
		Repeated:    repeated,
		Optional:    optional,
	}, nil
}

// parseEnum converts a EnumDescriptorProto to a schema.Enum.
func (p *Parser) parseEnum(enumDesc *descriptorpb.EnumDescriptorProto) *schema.Enum {
	enum := &schema.Enum{
		Name:   enumDesc.GetName(),
		Values: make([]string, 0),
	}

	for _, valueDesc := range enumDesc.Value {
		enum.Values = append(enum.Values, valueDesc.GetName())
	}

	return enum
}

// typeInfoFromProtoType converts a protobuf type to TypeInfo.
func (p *Parser) typeInfoFromProtoType(fieldType descriptorpb.FieldDescriptorProto_Type, typeName string) types.TypeInfo {
	ti := types.TypeInfo{
		Kind: types.KindScalar,
	}

	switch fieldType {
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		ti.Name = "string"
	case descriptorpb.FieldDescriptorProto_TYPE_INT32, descriptorpb.FieldDescriptorProto_TYPE_INT64:
		ti.Name = "int"
	case descriptorpb.FieldDescriptorProto_TYPE_UINT32, descriptorpb.FieldDescriptorProto_TYPE_UINT64:
		ti.Name = "uint"
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT, descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		ti.Name = "float64"
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		ti.Name = "bool"
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		ti.Kind = types.KindMessage
		// Remove package prefix from type name (e.g., ".pos.item_categories.v1.ItemCategory" -> "ItemCategory")
		if strings.HasPrefix(typeName, ".") {
			parts := strings.Split(typeName, ".")
			if len(parts) > 0 {
				ti.Name = parts[len(parts)-1]
			} else {
				ti.Name = typeName
			}
		} else {
			ti.Name = typeName
		}
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		ti.Kind = types.KindEnum
		// Remove package prefix from enum name
		if strings.HasPrefix(typeName, ".") {
			parts := strings.Split(typeName, ".")
			if len(parts) > 0 {
				ti.Name = parts[len(parts)-1]
			} else {
				ti.Name = typeName
			}
		} else {
			ti.Name = typeName
		}
	default:
		ti.Name = "string"
	}

	return ti
}

// mapProtoToDBType maps a protobuf TypeInfo to a database TypeInfo.
func (p *Parser) mapProtoToDBType(protoType types.TypeInfo) types.TypeInfo {
	// Simple mapping - in production this would use the resolver
	dbType := protoType
	return dbType
}

// isWellKnownType checks if a message name is a well-known protobuf type.
func isWellKnownType(name string) bool {
	wellKnownTypes := map[string]bool{
		"google.protobuf.Timestamp":   true,
		"google.protobuf.Duration":    true,
		"google.protobuf.DoubleValue": true,
		"google.protobuf.FloatValue":  true,
		"google.protobuf.Int64Value":  true,
		"google.protobuf.UInt64Value": true,
		"google.protobuf.Int32Value":  true,
		"google.protobuf.UInt32Value": true,
		"google.protobuf.BoolValue":   true,
		"google.protobuf.StringValue": true,
		"google.protobuf.BytesValue":  true,
		"google.protobuf.Any":         true,
		"google.protobuf.Empty":       true,
		"google.protobuf.Struct":      true,
		"google.protobuf.Value":       true,
		"google.protobuf.ListValue":   true,
		// Short names for well-known types
		"Timestamp":   true,
		"Duration":    true,
		"DoubleValue": true,
		"FloatValue":  true,
		"Int64Value":  true,
		"UInt64Value": true,
		"Int32Value":  true,
		"UInt32Value": true,
		"BoolValue":   true,
		"StringValue": true,
		"BytesValue":  true,
		"Any":         true,
		"Empty":       true,
		"Struct":      true,
		"Value":       true,
		"ListValue":   true,
	}
	return wellKnownTypes[name]
}

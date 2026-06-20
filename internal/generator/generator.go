package generator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jwart212/protoc-gen-go-mapper/internal/graph"
	"github.com/jwart212/protoc-gen-go-mapper/internal/handler"
	"github.com/jwart212/protoc-gen-go-mapper/internal/schema"
	"github.com/jwart212/protoc-gen-go-mapper/internal/template"
	"github.com/jwart212/protoc-gen-go-mapper/pkg/converter"
)

// Generator produces Go mapping code from Mapper graphs.
type Generator struct {
	tmpl         *template.Template
	generatedPkg string
	handlerReg   *handler.HandlerRegistry
}

// New creates a new Generator instance with loaded templates.
func New() *Generator {
	g := &Generator{
		tmpl:       template.New(),
		handlerReg: handler.NewHandlerRegistry(),
	}
	g.loadTemplates()
	return g
}

// SetHandlerRegistry sets the handler registry for field handling.
func (g *Generator) SetHandlerRegistry(reg *handler.HandlerRegistry) {
	g.handlerReg = reg
}

// SetGeneratedPackage sets the package name for the generated file.
func (g *Generator) SetGeneratedPackage(pkg string) {
	g.generatedPkg = pkg
}

// loadTemplates loads the code generation templates.
func (g *Generator) loadTemplates() {
	g.tmpl.Load("toProto", `func ToProto{{.Name}}(src db.{{.Name}}) *pb.{{.Name}} {
	return &pb.{{.Name}}{
		{{range $i, $f := .Fields}}{{if $i}}, {{end}}{{$f.Name}}: src.{{$f.Name}}{{end}}
	}
}
`)

	g.tmpl.Load("toDB", `func ToDB{{.Name}}(src *pb.{{.Name}}) db.{{.Name}} {
	return db.{{.Name}}{
		{{range $i, $f := .Fields}}{{if $i}}, {{end}}{{$f.Name}}: src.{{$f.Name}}{{end}}
	}
}
`)
}

// Generate produces mapping code for a single message.
func (g *Generator) Generate(msg *schema.Message, protoToDB, dbToProto *graph.Mapper, typeMappings map[string]string, useGenericConverters bool) (string, error) {
	var code string

	// Get the DB type name from type mappings, or use the proto message name
	dbTypeName := msg.Name
	if mappedName, ok := typeMappings[msg.Name]; ok {
		dbTypeName = mappedName
	}

	// Qualify the proto type name with its package if available and different from generated package
	// If the message package is the same as the generated package, don't qualify
	protoTypeName := msg.Name
	if msg.Package != "" && msg.Package != g.generatedPkg {
		protoTypeName = qualifyType(msg.Name, msg.Package)
	}

	// Sort fields by FieldNumber for deterministic output
	fields := make([]*schema.Field, len(msg.Fields))
	copy(fields, msg.Fields)
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].FieldNumber < fields[j].FieldNumber
	})

	// Generate ToProto function
	code += "func ToProto" + msg.Name + "(src sqlc." + dbTypeName + ") *" + protoTypeName + " {\n"
	code += "\treturn &" + protoTypeName + "{\n"
	for _, field := range fields {
		// Use PascalCase for DB fields, protoCase for proto fields
		dbFieldName := toPascalCase(field.Name)
		protoFieldName := toProtoCase(field.Name)

		// Check if a handler matches this field
		if g.handlerReg != nil {
			if h := g.handlerReg.Find(field, dbTypeName); h != nil {
				// Use handler for ToProto conversion
				expr, err := h.GenerateToProto(field, dbFieldName, protoFieldName, dbTypeName)
				if err != nil {
					return "", fmt.Errorf("handler error for field %s: %w", field.Name, err)
				}
				// If handler returns empty string, skip the field
				if expr == "" {
					continue
				}
				// Prepend src. prefix for ToProto conversions
				code += "\t\t" + strings.Replace(expr, ".", "src.", 1) + ",\n"
				continue
			}
		}

		// Find the converter for this field
		var conv converter.Converter
		for _, mapping := range dbToProto.Fields {
			if mapping.SourceField == field.Name {
				conv = mapping.Converter
				break
			}
		}
		if conv != nil && !useGenericConverters {
			expr, err := conv.Generate(converter.MappingField{
				SourceField: field.Name,
				TargetField: field.Name,
				SourceExpr:  "src." + dbFieldName,
				TargetExpr:  protoFieldName,
				SourceType:  field.DBType,
				TargetType:  field.ProtoType,
			})
			if err != nil {
				return "", err
			}
			code += "\t\t" + protoFieldName + ": " + expr + ",\n"
		} else if useGenericConverters {
			// Use generic converters if no converter found and generic mode is enabled
			protoFieldType := field.ProtoType.Name
			dbFieldType := field.DBType.Name

			// UUID fields (string <-> pgtype.UUID)
			if protoFieldType == "string" && dbFieldType == "pgtype.UUID" {
				if field.Optional {
					code += fmt.Sprintf("\t\t%s: ConvertUUID[*string](src.%s),\n", protoFieldName, dbFieldName)
				} else {
					code += fmt.Sprintf("\t\t%s: ConvertUUID[string](src.%s),\n", protoFieldName, dbFieldName)
				}
				// Timestamp fields (google.protobuf.Timestamp <-> pgtype.Timestamptz)
			} else if (protoFieldType == "Timestamp" || protoFieldType == "google.protobuf.Timestamp") && dbFieldType == "pgtype.Timestamptz" {
				code += fmt.Sprintf("\t\t%s: ConvertTimestamp[*timestamppb.Timestamp](src.%s),\n", protoFieldName, dbFieldName)
				// Text fields (string <-> pgtype.Text)
			} else if protoFieldType == "string" && dbFieldType == "pgtype.Text" {
				if field.Optional {
					code += fmt.Sprintf("\t\t%s: ConvertText[*string](src.%s),\n", protoFieldName, dbFieldName)
				} else {
					code += fmt.Sprintf("\t\t%s: ConvertText[string](src.%s),\n", protoFieldName, dbFieldName)
				}
				// Direct assignment for other types
			} else {
				code += "\t\t" + protoFieldName + ": src." + dbFieldName + ",\n"
			}
		} else {
			code += "\t\t" + protoFieldName + ": src." + dbFieldName + ",\n"
		}
	}
	code += "\t}\n}\n\n"

	// Generate ToDB function
	code += "func ToDB" + msg.Name + "(src *" + protoTypeName + ") sqlc." + dbTypeName + " {\n"
	code += "\treturn sqlc." + dbTypeName + "{\n"
	for _, field := range fields {
		// Use protoCase for proto fields, PascalCase for DB fields
		protoFieldName := toProtoCase(field.Name)
		dbFieldName := toPascalCase(field.Name)

		// Check if a handler matches this field
		if g.handlerReg != nil {
			if h := g.handlerReg.Find(field, dbTypeName); h != nil {
				// Use handler for ToDB conversion
				expr, err := h.GenerateToDB(field, dbFieldName, protoFieldName, dbTypeName)
				if err != nil {
					return "", fmt.Errorf("handler error for field %s: %w", field.Name, err)
				}
				// If handler returns empty string, skip the field
				if expr == "" {
					continue
				}
				code += "\t\t" + expr + ",\n"
				continue
			}
		}

		// Use generic converters if enabled (before registry lookup)
		if useGenericConverters {
			protoFieldType := field.ProtoType.Name
			dbFieldType := field.DBType.Name

			// UUID fields (string <-> pgtype.UUID)
			if protoFieldType == "string" && dbFieldType == "pgtype.UUID" {
				if field.Optional {
					code += fmt.Sprintf("\t\t%s: ConvertUUIDPtrToDB(src.%s),\n", dbFieldName, protoFieldName)
				} else {
					code += fmt.Sprintf("\t\t%s: ConvertUUIDToDB(src.%s),\n", dbFieldName, protoFieldName)
				}
				continue
				// Timestamp fields (google.protobuf.Timestamp <-> pgtype.Timestamptz)
			} else if (protoFieldType == "Timestamp" || protoFieldType == "google.protobuf.Timestamp") && dbFieldType == "pgtype.Timestamptz" {
				code += fmt.Sprintf("\t\t%s: ConvertTimestampToDB(src.%s),\n", dbFieldName, protoFieldName)
				continue
				// Text fields (string <-> pgtype.Text)
			} else if protoFieldType == "string" && dbFieldType == "pgtype.Text" {
				code += fmt.Sprintf("\t\t%s: ConvertTextToDB(src.%s),\n", dbFieldName, protoFieldName)
				continue
			}
		}

		// Find the converter for this field (fallback for non-generic mode or unhandled types)
		var conv converter.Converter
		for _, mapping := range protoToDB.Fields {
			if mapping.SourceField == field.Name {
				conv = mapping.Converter
				break
			}
		}
		if conv != nil && !useGenericConverters {
			expr, err := conv.Generate(converter.MappingField{
				SourceField: field.Name,
				TargetField: field.Name,
				SourceExpr:  "src." + protoFieldName,
				TargetExpr:  dbFieldName,
				SourceType:  field.ProtoType,
				TargetType:  field.DBType,
			})
			if err != nil {
				return "", err
			}
			code += "\t\t" + dbFieldName + ": " + expr + ",\n"
		} else {
			code += "\t\t" + dbFieldName + ": src." + protoFieldName + ",\n"
		}
	}
	code += "\t}\n}\n"

	return code, nil
}

// toPascalCase converts snake_case to PascalCase
func toPascalCase(s string) string {
	// Special case for "id" -> "ID"
	if s == "id" {
		return "ID"
	}
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			// Handle "id" suffix specially - keep it uppercase
			if part == "id" {
				parts[i] = "ID"
			} else {
				parts[i] = strings.ToUpper(string(part[0])) + strings.ToLower(part[1:])
			}
		}
	}
	return strings.Join(parts, "")
}

// toProtoCase converts snake_case to proto field case (first letter capitalized, rest camelCase)
func toProtoCase(s string) string {
	// Special case for "id" -> "Id"
	if s == "id" {
		return "Id"
	}
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			// Handle "id" suffix specially - keep it lowercase for proto
			if part == "id" {
				parts[i] = "Id"
			} else {
				if i == 0 {
					parts[i] = strings.ToUpper(string(part[0])) + strings.ToLower(part[1:])
				} else {
					parts[i] = strings.ToUpper(string(part[0])) + strings.ToLower(part[1:])
				}
			}
		}
	}
	return strings.Join(parts, "")
}

// toCamelCase converts snake_case to camelCase (for proto fields)
func toCamelCase(s string) string {
	// Special case for "id" -> "Id"
	if s == "id" {
		return "Id"
	}
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			if i == 0 {
				parts[i] = strings.ToLower(part)
			} else {
				parts[i] = strings.ToUpper(string(part[0])) + strings.ToLower(part[1:])
			}
		}
	}
	return strings.Join(parts, "")
}

// qualifyType qualifies a type name with its package prefix if the package is non-empty.
func qualifyType(name, pkg string) string {
	if pkg == "" {
		return name
	}
	return pkg + "." + name
}

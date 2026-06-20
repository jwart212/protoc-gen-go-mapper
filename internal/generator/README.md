# generator

Package generator produces Go mapping code from Mapper graphs.

## Overview

The generator package is the final stage of the compiler-style architecture. It consumes Mapper graphs and emits Go code for type conversions. The generator does not inspect proto descriptors - it only uses the Mapper graph.

## Types

### Generator

Generator produces Go mapping code from Mapper graphs:

```go
type Generator struct{}

func New() *Generator
func (g *Generator) Generate(msg *schema.Message, protoToDB, dbToProto *graph.Mapper) (string, error)
```

## Generated Functions

For each protobuf message, the generator produces:

```go
func ToProtoUser(src db.User) *pb.User
func ToDBUser(src *pb.User) db.User
```

## Determinism

The generator ensures deterministic output by:
- Sorting fields by FieldNumber before generation
- Never iterating over maps directly in output
- Producing byte-identical output for identical inputs

## Usage Example

```go
import "gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/internal/generator"

g := generator.New()
code, err := g.Generate(message, protoToDBMapper, dbToProtoMapper)
if err != nil {
    // Handle error
}
```

## Design Decisions

- **No descriptor inspection**: Generator only consumes Mapper graphs, never proto descriptors.
- **Deterministic output**: Fields are sorted by FieldNumber to ensure stable output.
- **Simple string generation**: Current implementation uses string concatenation (will be refactored to templates in M7).
- **Build-time validation**: All converter resolution happens during graph construction, not generation.

# proto

Package proto provides protobuf descriptor parsing into the internal schema model.

## Overview

The proto package implements the parser stage of the compiler-style architecture. It converts protobuf FileDescriptorProto objects into the internal schema representation.

## Types

### Parser

Parser converts protobuf descriptors into the internal schema model:

```go
type Parser struct{}

func New() *Parser
func (p *Parser) ParseFile() (*schema.Model, error)
```

## Usage Example

```go
import "github.com/jwart212/protoc-gen-go-mapper/internal/parser/proto"

p := proto.New()
model, err := p.ParseFile()
if err != nil {
    // Handle error
}
```

## Design Decisions

- **No type resolution**: The parser only converts descriptors to schema model without resolving types. Type resolution is handled by the resolver stage.
- **FieldNumber preservation**: The parser preserves FieldNumber from descriptors to ensure deterministic output.
- **Stateless**: Parser instances are stateless and can be reused.

## Future Implementation

The current ParseFile() is a placeholder. The full implementation will:
- Accept FileDescriptorProto as input
- Parse message definitions with their fields
- Parse enum definitions
- Preserve FieldNumber ordering from descriptors
- Return a complete schema.Model

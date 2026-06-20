# plugin

Package plugin orchestrates the complete mapping generation pipeline.

## Overview

The plugin package wires together all components into a coherent generation pipeline:
1. Load configuration
2. Parse proto files
3. Build mapper graphs
4. Generate code using templates

## Types

### Plugin

Plugin orchestrates the complete mapping generation pipeline:

```go
type Plugin struct {
    cfg       *config.Config
    parser    *proto.Parser
    registry  *registry.Registry
    generator *generator.Generator
}

func New(cfg *config.Config) *Plugin
func (p *Plugin) Generate(req *GenerateRequest, w io.Writer) error
```

### GenerateRequest

GenerateRequest represents a generation request:

```go
type GenerateRequest struct {
    ProtoFile string
}
```

## Pipeline

1. **Load Config**: Load and validate mapper.yaml
2. **Parse Proto**: Convert proto descriptors to schema model
3. **Register Converters**: Register all available converters
4. **Build Graphs**: Create mapper graphs for each message
5. **Generate Code**: Emit Go code using templates

## Usage Example

```go
import "gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/internal/plugin"

cfg, err := config.Load("mapper.yaml")
if err != nil {
    // Handle error
}

p := plugin.New(cfg)
req := &plugin.GenerateRequest{ProtoFile: "user.proto"}

err = p.Generate(req, os.Stdout)
if err != nil {
    // Handle error
}
```

## Design Decisions

- **Single entry point**: Plugin provides a single Generate() method for the entire pipeline.
- **Build-time validation**: All errors are caught during generation, not at runtime.
- **No descriptor inspection**: Generator only consumes Mapper graphs.

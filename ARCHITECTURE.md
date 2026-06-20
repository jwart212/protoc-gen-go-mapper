# protoc-gen-go-mapper Architecture Diagram

## High-Level Architecture

This diagram shows the overall flow from protoc invocation to generated code output. The Plugin orchestrates all components, coordinating configuration loading, proto parsing, type resolution, converter registration, field handler loading, and code generation.

**Components:**
- **protoc**: The Protocol Buffers compiler that invokes the plugin
- **Plugin**: Main orchestrator that coordinates all components
- **Config**: Loads and validates mapper.yaml configuration
- **Parser**: Parses proto files and extracts message definitions
- **Resolver**: Maps protobuf types to database-specific types
- **Registry**: Manages converter registration and resolution
- **HandlerRegistry**: Manages field handlers for special cases
- **Generator**: Produces Go mapping code
- **Generated Code**: Final output file with mapping functions

```mermaid
graph TD
    A[protoc] -->|Plugin Invocation| B[Plugin]
    B --> C[Config]
    B --> D[Parser]
    B --> E[Resolver]
    B --> F[Registry]
    B --> G[HandlerRegistry]
    B --> H[Generator]
    H --> I[Generated Code]
    
    C -->|mapper.yaml| J[Configuration]
    D -->|Proto File| K[Proto Descriptors]
    E -->|Type Mapping| L[Database Types]
    F -->|Converters| M[Converter Registry]
    G -->|Field Handlers| N[Handler Registry]
    
    style A fill:#e1f5ff
    style B fill:#fff4e1
    style I fill:#e8f5e9
```

## Component Interaction Diagram

This diagram illustrates the relationships between core and supporting components. The Plugin acts as the central coordinator, connecting to all core components (Parser, Resolver, Registry, Generator, HandlerRegistry) and the Config component. Supporting components (Schema, Graph, Template, Converter, Handler) provide type definitions, mapping structures, code templates, conversion logic, and field handling capabilities.

**Core Components:**
- **Plugin**: Central orchestrator
- **Parser**: Extracts message structures from proto files
- **Resolver**: Maps types between protobuf and database
- **Registry**: Manages converter registration
- **Generator**: Produces final Go code
- **HandlerRegistry**: Manages field-level custom logic

**Supporting Components:**
- **Config**: Configuration management
- **Schema**: Type system definitions
- **Graph**: Mapping graph structures
- **Template**: Code generation templates
- **Converter**: Type conversion implementations
- **Handler**: Field handler implementations

```mermaid
graph LR
    subgraph "Core Components"
        A[Plugin]
        B[Parser]
        C[Resolver]
        D[Registry]
        E[Generator]
        F[HandlerRegistry]
    end
    
    subgraph "Supporting Components"
        G[Config]
        H[Schema]
        I[Graph]
        J[Template]
        K[Converter]
        L[Handler]
    end
    
    A --> B
    A --> C
    A --> D
    A --> E
    A --> F
    A --> G
    
    B --> H
    C --> H
    D --> K
    E --> J
    F --> L
    
    E --> I
    
    style A fill:#fff4e1
    style B fill:#e1f5ff
    style C fill:#e1f5ff
    style D fill:#e1f5ff
    style E fill:#e1f5ff
    style F fill:#e1f5ff
```

## Data Flow Diagram

This sequence diagram shows the step-by-step execution flow from protoc invocation to file output. The process begins with protoc invoking the Plugin, which then loads configuration, parses proto files, resolves types, registers converters, loads field handlers, generates code, and finally writes the generated code to a file.

**Flow Steps:**
1. **protoc** invokes the Plugin
2. **Plugin** loads mapper.yaml configuration
3. **Config** returns validated configuration object
4. **Plugin** parses proto file to extract message definitions
5. **Parser** returns schema model with message structures
6. **Plugin** resolves protobuf types to database types
7. **Resolver** returns database type mappings
8. **Plugin** registers built-in and custom converters
9. **Registry** returns converter registry
10. **Plugin** loads field handlers from configuration
11. **HandlerRegistry** returns field handler registry
12. **Plugin** generates mapping code
13. **Generator** returns generated Go code
14. **Plugin** writes generated code to file

```mermaid
sequenceDiagram
    participant P as protoc
    participant PL as Plugin
    participant CFG as Config
    participant PR as Parser
    participant RES as Resolver
    participant REG as Registry
    participant HR as HandlerRegistry
    participant GEN as Generator
    participant OUT as Generated Code
    
    P->>PL: Invoke Plugin
    PL->>CFG: Load mapper.yaml
    CFG-->>PL: Config Object
    PL->>PR: Parse Proto File
    PR-->>PL: Schema Model
    PL->>RES: Resolve Types
    RES-->>PL: DB Type Mappings
    PL->>REG: Register Converters
    REG-->>PL: Converter Registry
    PL->>HR: Load Field Handlers
    HR-->>PL: Handler Registry
    PL->>GEN: Generate Code
    GEN-->>PL: Generated Go Code
    PL->>OUT: Write to File
```

## Database Resolver Architecture

This diagram shows how the Resolver maps protobuf types to database-specific types based on the configured database backend. The Resolver uses a switch statement to select the appropriate resolver (SQLC, PGX, or database_sql), each of which has its own type mapping strategy.

**Database Backends:**
- **SQLC**: Uses pgtype types for PostgreSQL (UUID, Timestamptz, Text, Numeric)
- **PGX**: Uses pgtype types with slight variations (UUID, Timestamp, Text, Numeric)
- **database_sql**: Uses standard Go types (string, time.Time, sql.NullString, sql.NullTime)

**Type Mappings:**
- **UUID**: pgtype.UUID (SQLC/PGX) or string (database_sql)
- **Timestamp**: pgtype.Timestamptz (SQLC), pgtype.Timestamp (PGX), or time.Time (database_sql)
- **Text**: pgtype.Text (SQLC/PGX) or string (database_sql)
- **Numeric**: pgtype.Numeric (SQLC/PGX) or string (database_sql)

```mermaid
graph TD
    A[Resolver] --> B{Database Type}
    
    B -->|sqlc| C[SQLC Resolver]
    B -->|pgx| D[PGX Resolver]
    B -->|database_sql| E[Database SQL Resolver]
    
    C --> F[pgtype.UUID]
    C --> G[pgtype.Timestamptz]
    C --> H[pgtype.Text]
    C --> I[pgtype.Numeric]
    
    D --> F
    D --> J[pgtype.Timestamp]
    D --> H
    D --> I
    
    E --> K[string]
    E --> L[time.Time]
    E --> M[sql.NullString]
    E --> N[sql.NullTime]
    
    style A fill:#fff4e1
    style C fill:#e1f5ff
    style D fill:#e1f5ff
    style E fill:#e1f5ff
```

## Handler System Architecture

This diagram shows the field handler system that provides flexible field-level customization. The HandlerRegistry manages multiple handler types, each with specific purposes. Handlers are resolved using priority-based matching, where higher priority handlers take precedence.

**Handler Types:**
- **type_assertion**: Handles type assertions for interface{} fields (e.g., converting interface{} to []string for SQLC array fields)
- **default_value**: Sets default values for fields that don't exist in the source (e.g., empty slice for tree children)
- **skip**: Skips fields during mapping (e.g., soft delete fields in responses)
- **field_mapping**: Provides custom expressions for both ToProto and ToDB directions

**Priority Levels:**
- **Priority 100**: Highest priority (e.g., skip handlers)
- **Priority 75**: Medium-high priority (e.g., default value handlers)
- **Priority 50**: Medium priority (e.g., type assertion handlers)

```mermaid
graph TD
    A[HandlerRegistry] --> B[Field Handlers]
    
    B --> C[type_assertion]
    B --> D[default_value]
    B --> E[skip]
    B --> F[field_mapping]
    
    C --> G[Interface to Type]
    D --> H[Set Default Value]
    E --> I[Skip Field]
    F --> J[Custom Expression]
    
    subgraph "Priority Resolution"
        K[Priority 100]
        L[Priority 75]
        M[Priority 50]
    end
    
    C --> M
    D --> L
    E --> K
    F --> M
    
    style A fill:#fff4e1
    style C fill:#e8f5e9
    style D fill:#e8f5e9
    style E fill:#fff3e0
    style F fill:#e8f5e9
```

## Configuration System

This diagram shows the structure of the mapper.yaml configuration file. The Config object loads and validates all configuration options, which control how the plugin generates mapping functions. Type conversions and field handlers support advanced features like pattern matching, priority-based resolution, and pointer strategies.

**Configuration Options:**
- **version**: Config version (must be "v1")
- **database**: Database type (sqlc, pgx, database_sql)
- **db_package**: Go package path for database models
- **package**: Proto and DB package names
- **type_mappings**: Custom proto message to DB model mappings
- **response_type_mappings**: Response message to SQLC Row type mappings
- **type_aliases**: Reusable type conversion definitions
- **type_conversions**: Custom type conversion rules with pattern matching
- **field_handlers**: Field-level custom logic
- **response_patterns**: Response helper configuration
- **pointer_settings**: Pointer handling strategies
- **messages**: List of messages to generate mappers for

**Advanced Features:**
- **Pattern Matching**: Regex-based field and message name matching
- **Priority**: Priority-based resolution for multiple matches
- **Pointer Strategy**: strict, lenient, or omit for nullable fields

```mermaid
graph TD
    A[mapper.yaml] --> B[Config]
    
    B --> C[version]
    B --> D[database]
    B --> E[db_package]
    B --> F[package]
    B --> G[type_mappings]
    B --> H[response_type_mappings]
    B --> I[type_aliases]
    B --> J[type_conversions]
    B --> K[field_handlers]
    B --> L[response_patterns]
    B --> M[pointer_settings]
    B --> N[messages]
    
    J --> O[Pattern Matching]
    J --> P[Priority]
    J --> Q[Pointer Strategy]
    
    K --> R[Field Name]
    K --> S[DB Type]
    K --> T[Message]
    K --> U[Priority]
    
    style A fill:#e1f5ff
    style B fill:#fff4e1
    style J fill:#e8f5e9
    style K fill:#e8f5e9
```

## Converter Registry

This diagram shows the converter registry system that manages type conversion logic. The Registry maintains a collection of converters for different type pairs (Scalar, UUID, Timestamp, Decimal, Enum, Nullable, Slice, Message). When type_conversions is empty (zero-config mode), generic converters (ConvertUUID, ConvertTimestamp, ConvertText) are automatically used. The registry uses priority-based resolution to select the best converter when multiple converters match a type pair.

**Built-in Converters:**
- **ScalarConverter**: Handles basic scalar types (int32, int64, string, bool, float64)
- **UUIDConverter**: Handles UUID type conversions
- **TimestampConverter**: Handles timestamp type conversions
- **DecimalConverter**: Handles decimal/numeric type conversions
- **EnumConverter**: Handles enum type conversions
- **NullableConverter**: Handles nullable/optional type conversions
- **SliceConverter**: Handles array/slice type conversions
- **MessageConverter**: Handles nested message type conversions

**Generic Converters (Zero-Config Mode):**
- **ConvertUUID**: Automatic UUID ↔ string conversion
- **ConvertTimestamp**: Automatic Timestamp ↔ time.Time conversion
- **ConvertText**: Automatic Text ↔ string conversion

**Resolution Process:**
1. Match type pair against all registered converters
2. Get priority for each matching converter
3. Select converter with highest priority
4. Return error if multiple converters have equal priority (ambiguous mapping)

```mermaid
graph TD
    A[Registry] --> B[Converters]
    
    B --> C[ScalarConverter]
    B --> D[UUIDConverter]
    B --> E[TimestampConverter]
    B --> F[DecimalConverter]
    B --> G[EnumConverter]
    B --> H[NullableConverter]
    B --> I[SliceConverter]
    B --> J[MessageConverter]
    
    subgraph "Generic Converters"
        K[ConvertUUID]
        L[ConvertTimestamp]
        M[ConvertText]
    end
    
    B --> K
    B --> L
    B --> M
    
    subgraph "Priority Resolution"
        N[Match Type Pair]
        O[Get Priority]
        P[Select Best]
    end
    
    C --> N
    D --> N
    E --> N
    N --> O
    O --> P
    
    style A fill:#fff4e1
    style K fill:#e8f5e9
    style L fill:#e8f5e9
    style M fill:#e8f5e9
```

## Package Structure

This diagram shows the overall package organization of the protoc-gen-go-mapper project. The project is divided into four main directories: cmd (command-line interface), internal (core implementation), pkg (public packages), and examples (sample implementations).

**Directory Structure:**
- **cmd**: Contains the main protoc-gen-go-mapper command-line tool
- **internal**: Core implementation packages (config, converter, generator, graph, handler, parser, registry, resolver, schema, template, plugin)
- **pkg**: Public packages (converter, errors, naming, types)
- **examples**: Sample implementations (simple, medium, complex, advanced)

**Internal Packages:**
- **config**: Configuration loading and validation
- **converter**: Generic converter implementations
- **generator**: Code generation logic
- **graph**: Mapping graph structures
- **handler**: Field handler implementations
- **parser**: Proto file parsing
- **registry**: Converter registration and resolution
- **resolver**: Type resolution for different databases
- **schema**: Type system definitions
- **template**: Code generation templates
- **plugin**: Main plugin orchestration

**Public Packages:**
- **converter**: Public converter interfaces
- **errors**: Error definitions
- **naming**: Naming conventions utilities
- **types**: Type system utilities

```mermaid
graph TD
    A[protoc-gen-go-mapper] --> B[cmd]
    A --> C[internal]
    A --> D[pkg]
    A --> E[examples]
    
    B --> F[protoc-gen-go-mapper]
    
    C --> G[config]
    C --> H[converter]
    C --> I[generator]
    C --> J[graph]
    C --> K[handler]
    C --> L[parser]
    C --> M[registry]
    C --> N[resolver]
    C --> O[schema]
    C --> P[template]
    C --> Q[plugin]
    
    D --> R[converter]
    D --> S[errors]
    D --> T[naming]
    D --> U[types]
    
    E --> V[simple]
    E --> W[Medium]
    E --> X[Complex]
    E --> Y[Advanced]
    
    style A fill:#e1f5ff
    style C fill:#fff4e1
    style D fill:#e8f5e9
    style E fill:#f3e5f5
```

## Zero-Config Mode Flow

This diagram shows the zero-config mode flow, which is activated when type_conversions is empty in the configuration. In this mode, the plugin automatically uses built-in generic converters for common types (UUID, Timestamp, Text) without requiring explicit configuration. The generic converters are inlined directly into the generated code, making it self-contained without external dependencies.

**Zero-Config Process:**
1. Check if type_conversions is empty
2. If empty, enable generic converters
3. Auto-detect field types (UUID, Timestamp, Text)
4. Apply appropriate generic converter
5. Inline converter functions in generated code
6. Generate final mapping code

**Generic Converters:**
- **ConvertUUID**: Handles UUID ↔ string conversions for both nullable and non-nullable fields
- **ConvertTimestamp**: Handles Timestamp ↔ time.Time conversions
- **ConvertText**: Handles Text ↔ string conversions for both nullable and non-nullable fields

**Benefits:**
- No configuration required for common types
- Self-contained generated code
- Automatic type detection
- Reduced configuration complexity

```mermaid
graph TD
    A[Start] --> B{type_conversions empty?}
    B -->|Yes| C[Enable Generic Converters]
    B -->|No| D[Use Custom Conversions]
    
    C --> E[Auto-detect Types]
    E --> F[UUID Fields]
    E --> G[Timestamp Fields]
    E --> H[Text Fields]
    
    F --> I[ConvertUUID]
    G --> J[ConvertTimestamp]
    H --> K[ConvertText]
    
    I --> L[Inline in Generated Code]
    J --> L
    K --> L
    
    D --> M[Use Registry]
    M --> N[Apply Custom Logic]
    
    L --> O[Generate Code]
    N --> O
    
    style C fill:#e8f5e9
    style D fill:#fff3e0
    style L fill:#e1f5ff
```

## Complete Pipeline

This diagram shows the complete end-to-end pipeline from user invocation to generated code output. The process begins with the user running a protoc command, which invokes the Plugin. The Plugin then loads and validates configuration, parses proto files, resolves types, registers converters, loads field handlers, generates mapping functions (ToProto, ToDB, and Response Helpers), and finally writes the generated code to a .proto_mapper.pb.go file.

**Pipeline Stages:**
1. **User Invocation**: User runs protoc with the mapper plugin
2. **Plugin Initialization**: Plugin.New creates plugin instance
3. **Configuration Loading**: Config.Load reads mapper.yaml
4. **Configuration Validation**: Config.Validate checks configuration
5. **Generation**: Plugin.Generate orchestrates the generation process
6. **Proto Parsing**: Parser.ParseFile extracts message definitions
7. **Type Resolution**: Resolver.Resolve maps types to database types
8. **Converter Registration**: Registry.Register registers converters
9. **Handler Loading**: HandlerRegistry.Load loads field handlers
10. **Code Generation**: Generator.Generate produces mapping functions
11. **Output**: Generated code written to .proto_mapper.pb.go file

**Generated Functions:**
- **ToProto Functions**: Convert database models to protobuf messages
- **ToDB Functions**: Convert protobuf messages to database models
- **Response Helpers**: Specialized functions for list responses

```mermaid
graph TD
    A[User] -->|protoc command| B[protoc]
    B -->|Plugin| C[Plugin.New]
    C -->|Load Config| D[Config.Load]
    D -->|Validate| E[Config.Validate]
    E -->|Valid| F[Plugin.Generate]
    
    F --> G[Parser.ParseFile]
    G --> H[Schema Model]
    
    F --> I[Resolver.Resolve]
    I --> J[DB Type Mappings]
    
    F --> K[Registry.Register]
    K --> L[Converter Registry]
    
    F --> M[HandlerRegistry.Load]
    M --> N[Field Handlers]
    
    F --> O[Generator.Generate]
    O --> P[ToProto Functions]
    O --> Q[ToDB Functions]
    O --> R[Response Helpers]
    
    P --> S[Generated Code]
    Q --> S
    R --> S
    
    S --> T[Write to File]
    T --> U[.proto_mapper.pb.go]
    
    style A fill:#e1f5ff
    style C fill:#fff4e1
    style O fill:#e8f5e9
    style U fill:#f3e5f5
```

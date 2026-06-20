# naming

Package naming provides name transformation utilities for converting between protobuf and database naming conventions.

## Overview

The naming package handles conversion between different naming conventions:
- snake_case ↔ camelCase ↔ PascalCase
- Database table names ↔ Protobuf message names

## Functions

### ToCamelCase

Converts a snake_case or SCREAMING_SNAKE_CASE field name to camelCase.

```go
ToCamelCase("user_name") // "userName"
ToCamelCase("USER_NAME") // "userName"
```

### ToPascalCase

Converts a snake_case or SCREAMING_SNAKE_CASE field name to PascalCase.

```go
ToPascalCase("user_name") // "UserName"
ToPascalCase("USER_NAME") // "UserName"
```

### ToProtoMessageName

Converts a database table name to a protobuf message name (snake_case → PascalCase).

```go
ToProtoMessageName("users") // "Users"
ToProtoMessageName("user_profiles") // "UserProfiles"
```

### ToDBTableName

Converts a protobuf message name to a database table name (PascalCase → snake_case).

```go
ToDBTableName("UserProfile") // "user_profile"
ToDBTableName("User") // "user"
```

## Usage Example

```go
import "gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/naming"

// Convert database column to protobuf field
protoField := naming.ToCamelCase("user_id") // "userId"

// Convert protobuf message to database table
dbTable := naming.ToDBTableName("UserProfile") // "user_profile"
```

## Design Decisions

- **Simple transformations**: Functions handle common cases without over-engineering for edge cases.
- **Deterministic output**: All transformations produce consistent, predictable results.
- **No external dependencies**: Pure Go implementation without regex or external libraries.

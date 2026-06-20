# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial implementation of protoc-gen-go-mapper
- Compiler-style architecture with parser, registry, graph, and generator
- Support for scalar type conversions
- UUID, Timestamp, Decimal, Enum, Nullable, and Slice converters
- Template-based code generation
- Configuration via mapper.yaml

### TODO
- Implement protobuf descriptor parsing
- Implement protoc plugin protocol (stdin/stdout)
- Add resolver package for database-specific type mapping
- Add integration tests with actual .proto files
- Add examples and documentation

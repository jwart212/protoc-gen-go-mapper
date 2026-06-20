# Contributing to protoc-gen-go-mapper

Thank you for your interest in contributing to protoc-gen-go-mapper!

## Development Setup

1. Clone the repository:
```bash
git clone https://github.com/jwart212/protoc-gen-go-mapper.git
cd protoc-gen-go-mapper
```

2. Install dependencies:
```bash
go mod download
```

3. Run tests:
```bash
go test ./...
```

## Code Style

This project follows:
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Project-specific rules in `.windsurf/rules/rules.md`

## Submitting Changes

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Ensure all tests pass
5. Submit a pull request

## Testing

- Unit tests for each package
- Golden tests for code generation
- Integration tests with actual .proto files

## Architecture

The project follows a compiler-style architecture:
Parser → Schema Model → Resolver → Converter Registry → Mapper Graph → Generator

See the main README.md for details.

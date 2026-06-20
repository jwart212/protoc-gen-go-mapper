# Medium Example

A medium example demonstrating protoc-gen-go-mapper with multiple tables and relationships.

## Schema

PostgreSQL schema with three related tables:
- `authors`: Authors with bio information
- `books`: Books linked to authors (foreign key)
- `reviews`: Reviews linked to books (foreign key)

## Relationships

- One author → Many books (author_id in books)
- One book → Many reviews (book_id in reviews)

## Setup

1. Generate sqlc code:
```bash
cd examples/medium
sqlc generate
```

2. Generate protobuf code:
```bash
protoc --go_out=. --go_opt=paths=source_relative library.proto
```

3. Generate mapper code:
```bash
protoc --go-mapper_out=. --go-mapper_opt=mapper.example.yaml library.proto
```

## Generated Functions

For each message, the plugin generates:
- `ToProtoAuthor(db.Author) *gen.Author`
- `ToDBAuthor(*gen.Author) db.Author`
- `ToProtoBook(db.Book) *gen.Book`
- `ToDBBook(*gen.Book) db.Book`
- `ToProtoReview(db.Review) *gen.Review`
- `ToDBReview(*gen.Review) db.Review`

## Usage

```go
// Convert DB author to protobuf
protoAuthor := ToProtoAuthor(dbAuthor)

// Convert protobuf book to DB
dbBook := ToDBBook(protoBook)
```

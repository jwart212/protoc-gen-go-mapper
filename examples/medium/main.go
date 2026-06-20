package main

import (
	"fmt"

	"github.com/jwart212/protoc-gen-go-mapper/examples/medium/gen"
)

func main() {
	// Create a proto author
	bio := "Software Engineer"
	protoAuthor := gen.Author{
		Id:    1,
		Name:  "Jane Doe",
		Email: "jane@example.com",
		Bio:   &bio,
	}

	// Convert Proto author to DB author
	dbAuthor := gen.ToDBAuthor(&protoAuthor)
	fmt.Printf("DB Author: %+v\n", dbAuthor)

	// Convert DB author back to Proto author
	protoAuthor2 := gen.ToProtoAuthor(dbAuthor)
	fmt.Printf("Proto Author: %+v\n", protoAuthor2)

	// Create a proto book
	publishedDate := "2024-01-01"
	authorId := int32(1)
	protoBook := gen.Book{
		Id:            1,
		Title:         "Go Programming",
		Isbn:          "123-4567890123",
		PublishedDate: &publishedDate,
		AuthorId:      &authorId,
	}

	// Convert Proto book to DB book
	dbBook := gen.ToDBBook(&protoBook)
	fmt.Printf("DB Book: %+v\n", dbBook)

	// Convert DB book back to Proto book
	protoBook2 := gen.ToProtoBook(dbBook)
	fmt.Printf("Proto Book: %+v\n", protoBook2)

	fmt.Println("Medium example completed successfully!")
}

package proto_test

import (
	"fmt"

	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/internal/parser/proto"
)

func ExampleNew() {
	p := proto.New("sqlc")
	fmt.Printf("Parser created: %v", p != nil)
	// Output: Parser created: true
}

package proto_test

import (
	"fmt"

	"github.com/jwart212/protoc-gen-go-mapper/internal/parser/proto"
)

func ExampleNew() {
	p := proto.New("sqlc")
	fmt.Printf("Parser created: %v", p != nil)
	// Output: Parser created: true
}

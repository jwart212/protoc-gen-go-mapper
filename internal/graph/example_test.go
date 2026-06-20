package graph_test

import (
	"fmt"

	"github.com/jwart212/protoc-gen-go-mapper/internal/graph"
	"github.com/jwart212/protoc-gen-go-mapper/internal/registry"
	"github.com/jwart212/protoc-gen-go-mapper/pkg/types"
)

func ExampleNewMapper() {
	m := graph.NewMapper("User", "User")
	fmt.Printf("Mapper: %s -> %s", m.Source, m.Target)
	// Output: Mapper: User -> User
}

func ExampleMapper_AddField() {
	m := graph.NewMapper("User", "User")
	r := registry.New()
	r.Register(registry.ScalarConverter{})

	srcType := types.TypeInfo{Kind: types.KindScalar}
	dstType := types.TypeInfo{Kind: types.KindScalar}

	m.AddField("id", "id", srcType, dstType, r)
	fmt.Printf("Fields: %d", len(m.Fields))
	// Output: Fields: 1
}

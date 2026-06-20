package types_test

import (
	"fmt"

	"github.com/jwart212/protoc-gen-go-mapper/pkg/types"
)

func ExampleTypeInfo() {
	ti := types.TypeInfo{
		Package: "github.com/google/uuid",
		Name:    "UUID",
		Kind:    types.KindUUID,
	}
	fmt.Printf("Type: %s.%s (Kind: %s)", ti.Package, ti.Name, ti.Kind)
	// Output: Type: github.com/google/uuid.UUID (Kind: UUID)
}

func ExampleKind_String() {
	fmt.Println(types.KindUUID.String())
	// Output: UUID
}

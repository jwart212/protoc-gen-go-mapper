package registry_test

import (
	"fmt"

	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/internal/registry"
	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/types"
)

func ExampleNew() {
	r := registry.New()
	fmt.Printf("Registry created: %v", r != nil)
	// Output: Registry created: true
}

func ExampleRegistry_Register() {
	r := registry.New()
	r.Register(registry.ScalarConverter{})
	fmt.Println("Converter registered")
	// Output: Converter registered
}

func ExampleRegistry_Resolve() {
	r := registry.New()
	r.Register(registry.ScalarConverter{})

	src := types.TypeInfo{Kind: types.KindScalar}
	dst := types.TypeInfo{Kind: types.KindScalar}

	c, err := r.Resolve(src, dst)
	if err == nil {
		fmt.Printf("Converter resolved with priority %d", c.Priority())
	}
	// Output: Converter resolved with priority 0
}

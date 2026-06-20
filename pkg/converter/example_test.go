package converter_test

import (
	"fmt"

	"github.com/jwart212/protoc-gen-go-mapper/pkg/converter"
	"github.com/jwart212/protoc-gen-go-mapper/pkg/types"
)

type exampleConverter struct{}

func (c exampleConverter) Match(src, dst types.TypeInfo) bool {
	return src.Kind == types.KindScalar && dst.Kind == types.KindScalar
}

func (c exampleConverter) Priority() int {
	return 0
}

func (c exampleConverter) Generate(field converter.MappingField) (string, error) {
	return field.SourceExpr, nil
}

func ExampleConverter() {
	c := exampleConverter{}
	fmt.Printf("Priority: %d", c.Priority())
	// Output: Priority: 0
}

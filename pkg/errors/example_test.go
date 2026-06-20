package errors_test

import (
	"errors"
	"fmt"

	mappererrors "github.com/jwart212/protoc-gen-go-mapper/pkg/errors"
)

func ExampleErrNoConverterFound() {
	err := mappererrors.ErrNoConverterFound
	fmt.Println(err)
	// Output: mapper: no converter found for type pair
}

func ExampleErrNoConverterFound_wrapping() {
	wrapped := fmt.Errorf("context: %w", mappererrors.ErrNoConverterFound)
	if errors.Is(wrapped, mappererrors.ErrNoConverterFound) {
		fmt.Println("Detected ErrNoConverterFound")
	}
	// Output: Detected ErrNoConverterFound
}

package naming_test

import (
	"fmt"

	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/pkg/naming"
)

func ExampleToCamelCase() {
	fmt.Println(naming.ToCamelCase("user_name"))
	// Output: userName
}

func ExampleToPascalCase() {
	fmt.Println(naming.ToPascalCase("user_name"))
	// Output: UserName
}

func ExampleToProtoMessageName() {
	fmt.Println(naming.ToProtoMessageName("users"))
	// Output: Users
}

func ExampleToDBTableName() {
	fmt.Println(naming.ToDBTableName("UserProfile"))
	// Output: user_profile
}

package main

import (
	"database/sql"
	"fmt"

	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/examples/simple/db"
	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/examples/simple/gen"
)

func main() {
	// Create a DB user
	dbUser := db.User{
		ID:     1,
		Name:   "John Doe",
		Email:  "john@example.com",
		Age:    sql.NullInt32{Int32: 30, Valid: true},
		Active: sql.NullBool{Bool: true, Valid: true},
	}

	// Convert DB user to Proto user
	protoUser := gen.ToProtoUser(dbUser)
	fmt.Printf("Proto User: %+v\n", protoUser)

	// Convert Proto user back to DB user
	dbUser2 := gen.ToDBUser(protoUser)
	fmt.Printf("DB User: %+v\n", dbUser2)

	fmt.Println("Simple example completed successfully!")
}

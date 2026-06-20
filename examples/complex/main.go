package main

import (
	"fmt"
	"time"

	"github.com/jwart212/protoc-gen-go-mapper/examples/complex/gen"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	// Create a proto customer
	phone := "+1234567890"
	createdAt := timestamppb.New(time.Now())
	updatedAt := timestamppb.New(time.Now())
	protoCustomer := gen.Customer{
		Id:        "550e8400-e29b-41d4-a716-446655440000",
		Name:      "John Doe",
		Email:     "john@example.com",
		Phone:     &phone,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	// Convert Proto customer to DB customer
	dbCustomer := gen.ToDBCustomer(&protoCustomer)
	fmt.Printf("DB Customer: %+v\n", dbCustomer)

	// Convert DB customer back to Proto customer
	protoCustomer2 := gen.ToProtoCustomer(dbCustomer)
	fmt.Printf("Proto Customer: %+v\n", protoCustomer2)

	// Create a proto order
	notes := "Test order"
	protoOrder := gen.Order{
		Id:          "550e8400-e29b-41d4-a716-446655440001",
		TotalAmount: "100.50",
		Notes:       &notes,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	// Convert Proto order to DB order
	dbOrder := gen.ToDBOrder(&protoOrder)
	fmt.Printf("DB Order: %+v\n", dbOrder)

	// Convert DB order back to Proto order
	protoOrder2 := gen.ToProtoOrder(dbOrder)
	fmt.Printf("Proto Order: %+v\n", protoOrder2)

	fmt.Println("Complex example completed successfully!")
}

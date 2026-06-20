# Complex Example

A complex example demonstrating protoc-gen-go-mapper with nested messages, enums, nullable types, UUID, and timestamps.

## Schema

PostgreSQL schema with advanced features:
- **UUID** primary keys (using gen_random_uuid())
- **Custom enum** (order_status)
- **Nullable fields** (phone, notes, description)
- **Timestamps** (created_at, updated_at)
- **Decimal types** (price, total_amount)
- **Foreign key relationships** (customers → orders → order_items → products)

## Tables

- `customers`: Customer information with UUID primary key
- `products`: Product catalog with pricing
- `orders`: Order management with status enum
- `order_items`: Order line items with product references

## Advanced Features Demonstrated

1. **UUID ↔ string conversion**: UUID fields mapped to string in protobuf
2. **Enum ↔ string conversion**: Custom PostgreSQL enum to protobuf enum
3. **Timestamp ↔ string conversion**: Timestamp fields mapped to string
4. **Decimal ↔ string conversion**: Decimal/numeric types mapped to string
5. **Nullable types**: Optional fields (phone, notes, description)
6. **Nested messages**: Order contains repeated OrderItem
7. **Repeated fields**: Order.items as repeated OrderItem

## Setup

1. Generate sqlc code:
```bash
cd examples/complex
sqlc generate
```

2. Generate protobuf code:
```bash
protoc --go_out=. --go_opt=paths=source_relative ecommerce.proto
```

3. Generate mapper code:
```bash
protoc --go-mapper_out=. --go-mapper_opt=mapper.example.yaml ecommerce.proto
```

## Generated Functions

For each message, the plugin generates:
- `ToProtoCustomer(db.Customer) *gen.Customer`
- `ToDBCustomer(*gen.Customer) db.Customer`
- `ToProtoProduct(db.Product) *gen.Product`
- `ToDBProduct(*gen.Product) db.Product`
- `ToProtoOrderItem(db.OrderItem) *gen.OrderItem`
- `ToDBOrderItem(*gen.OrderItem) db.OrderItem`
- `ToProtoOrder(db.Order) *gen.Order`
- `ToDBOrder(*gen.Order) db.Order`

## Type Conversions

| PostgreSQL | Protobuf | Converter |
|------------|----------|-----------|
| UUID | string | UUIDConverter |
| order_status (enum) | OrderStatus (enum) | EnumConverter |
| TIMESTAMP | string | TimestampConverter |
| DECIMAL | string | DecimalConverter |
| TEXT (nullable) | string (optional) | NullableConverter |
| INTEGER | int32 | ScalarConverter |

## Usage

```go
// Convert DB customer to protobuf
protoCustomer := ToProtoCustomer(dbCustomer)

// Convert protobuf order to DB
dbOrder := ToDBOrder(protoOrder)

// Handle nested items
for _, item := range protoOrder.Items {
    dbItem := ToDBOrderItem(item)
}
```

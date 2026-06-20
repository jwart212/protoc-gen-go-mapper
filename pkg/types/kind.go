package types

// Kind classifies a TypeInfo for converter matching.
type Kind int

const (
	KindScalar Kind = iota
	KindUUID
	KindTimestamp
	KindDecimal
	KindEnum
	KindNullable
	KindMessage
)

// String returns a human-readable representation of Kind.
func (k Kind) String() string {
	switch k {
	case KindScalar:
		return "Scalar"
	case KindUUID:
		return "UUID"
	case KindTimestamp:
		return "Timestamp"
	case KindDecimal:
		return "Decimal"
	case KindEnum:
		return "Enum"
	case KindNullable:
		return "Nullable"
	case KindMessage:
		return "Message"
	default:
		return "Unknown"
	}
}

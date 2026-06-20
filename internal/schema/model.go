package schema

// Model represents the complete schema parsed from protobuf descriptors.
type Model struct {
	Messages []*Message
	Enums    []*Enum
}

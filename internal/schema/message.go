package schema

// Message represents a protobuf message definition.
type Message struct {
	Name    string
	Package string
	Fields  []*Field
}

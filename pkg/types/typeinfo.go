package types

// TypeInfo describes a single Go type as understood by the mapper.
// It deliberately avoids raw strings outside the registry package.
type TypeInfo struct {
	Package string
	Name    string

	IsPointer  bool
	IsSlice    bool
	IsEnum     bool
	IsNullable bool

	Kind Kind
}

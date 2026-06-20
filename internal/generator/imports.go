package generator

import (
	"fmt"
	"sort"
	"strings"
)

// Imports manages import statements for generated code.
type Imports struct {
	standard   map[string]bool
	thirdParty map[string]bool
	local      map[string]bool
}

// NewImports creates a new Imports instance.
func NewImports() *Imports {
	return &Imports{
		standard:   make(map[string]bool),
		thirdParty: make(map[string]bool),
		local:      make(map[string]bool),
	}
}

// AddStandard adds a standard library import.
func (i *Imports) AddStandard(pkg string) {
	i.standard[pkg] = true
}

// AddThirdParty adds a third-party import.
func (i *Imports) AddThirdParty(pkg string) {
	i.thirdParty[pkg] = true
}

// AddLocal adds a local import.
func (i *Imports) AddLocal(pkg string) {
	i.local[pkg] = true
}

// Generate generates the import block.
func (i *Imports) Generate() string {
	var imports []string

	// Standard library
	for pkg := range i.standard {
		imports = append(imports, fmt.Sprintf(`"%s"`, pkg))
	}

	// Third-party
	for pkg := range i.thirdParty {
		imports = append(imports, fmt.Sprintf(`"%s"`, pkg))
	}

	// Local
	for pkg := range i.local {
		imports = append(imports, fmt.Sprintf(`"%s"`, pkg))
	}

	sort.Strings(imports)

	if len(imports) == 0 {
		return ""
	}

	return fmt.Sprintf("import (\n\t%s\n)", strings.Join(imports, "\n\t"))
}

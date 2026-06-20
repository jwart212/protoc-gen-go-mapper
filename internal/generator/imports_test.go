package generator

import (
	"testing"
)

func TestNewImports(t *testing.T) {
	i := NewImports()
	if i == nil {
		t.Error("NewImports() returned nil")
	}
}

func TestImports(t *testing.T) {
	i := NewImports()

	i.AddStandard("fmt")
	i.AddThirdParty("github.com/example/pkg")
	i.AddLocal("internal/gen")

	if !i.standard["fmt"] {
		t.Error("Expected fmt in standard imports")
	}
	if !i.thirdParty["github.com/example/pkg"] {
		t.Error("Expected github.com/example/pkg in third-party imports")
	}
	if !i.local["internal/gen"] {
		t.Error("Expected internal/gen in local imports")
	}
}

func TestImportsGenerate(t *testing.T) {
	i := NewImports()

	i.AddStandard("fmt")
	i.AddThirdParty("github.com/example/pkg")

	got := i.Generate()
	if got == "" {
		t.Error("Generate() returned empty string")
	}
}

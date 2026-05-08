package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/ast"
	"github.com/gonzalomdvc/go-linter/packages"
	linttypes "github.com/gonzalomdvc/go-linter/types"
)

func Test_GL14(t *testing.T) {
	state := &packages.State{
		SourceAsts: make(map[string]packages.SourceAst),
	}
	// Load state with ASTs for test files
	astFile, fset, err := ast.GetAst("../test/GL14.go", nil)
	if err != nil {
		t.Fatalf("Expected no error generating AST for GL14.go, got: %v", err)
	}
	state.SourceAsts["../test/GL14.go"] = packages.SourceAst{
		Fset:    fset,
		AstFile: astFile,
	}
	err = linttypes.ResolveTypesInfoForFiles(state, fset)
	if err != nil {
		t.Fatalf("Expected no error preloading types info, got %v", err)
	}

	positions := []Position{
		{
			Line:   4,
			Column: 2,
		},
		{
			Line:   5,
			Column: 2,
		},
		{
			Line:   12,
			Column: 2,
		},
		{
			Line:   13,
			Column: 2,
		},
	}
	err = RunCheckTest("GL14.go", true, positions, GL14, state)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

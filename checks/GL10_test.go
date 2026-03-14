package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/ast"
	"github.com/gonzalomdvc/go-linter/packages"
)

func Test_GL10(t *testing.T) {
	positions := []Position{
		{
			Column: 9,
			Line:   8,
		},
	}
	state := &packages.State{
		Packages: make(map[string]packages.Package),
	}
	funcDeclsCh := make(chan packages.FuncDeclResult, 1)

	astFile, _, err := ast.GetAst("../test/GL10.go")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	packages.ImportPackages(astFile, funcDeclsCh, state)
	funcDeclResult := <-funcDeclsCh

	if _, exists := state.Packages[funcDeclResult.PackagePath]; !exists {
		state.Packages[funcDeclResult.PackagePath] = packages.Package{FuncDecls: funcDeclResult.FuncDecls}
	}

	err = RunCheckTest("GL10.go", true, positions, GL10, state)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

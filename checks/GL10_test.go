package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/ast"
	"github.com/gonzalomdvc/go-linter/interfaces"
	"github.com/gonzalomdvc/go-linter/packages"
	"github.com/gonzalomdvc/go-linter/test"
)

func Test_GL10(t *testing.T) {
	positions := []interfaces.Position{
		{
			Column: 9,
			Line:   8,
		},
	}
	state := &interfaces.State{
		Packages: make(map[string]interfaces.Package),
	}
	funcDeclsCh := make(chan packages.FuncDeclResult, 1)

	go func() {
		for funcDeclResult := range funcDeclsCh {
			if _, exists := state.Packages[funcDeclResult.PackagePath]; !exists {
				state.Packages[funcDeclResult.PackagePath] = interfaces.Package{FuncDecls: funcDeclResult.FuncDecls}
			}
		}
	}()

	astFile, _, err := ast.GetAst("../test/GL10.go")
	packages.ImportPackages(astFile, funcDeclsCh, state)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	err = test.RunCheckTest("GL10.go", true, positions, GL10, state)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

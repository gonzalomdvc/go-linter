package types

import (
	"go/token"
	"path/filepath"
	"testing"

	"github.com/gonzalomdvc/go-linter/ast"
	"github.com/gonzalomdvc/go-linter/packages"
)

func TestResolveTypesInfoForFiles(t *testing.T) {
	state := &packages.State{
		SourceAsts: make(map[string]packages.SourceAst),
	}
	fset := token.NewFileSet()
	// Pre-populate state with ASTs for test files
	testFiles := []string{
		"../test/GL10_helper/GL10_helper.go",
		"../test/GL14.go",
	}
	for _, filePath := range testFiles {
		astFile, fset, err := ast.GetAst(filePath, fset)
		if err != nil {
			t.Fatalf("Expected no error generating AST for %s, got: %v", filePath, err)
		}
		state.SourceAsts[filePath] = packages.SourceAst{
			Fset:    fset,
			AstFile: astFile,
		}
	}
	err := ResolveTypesInfoForFiles(state, fset)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(state.TypesInfo) != 2 {
		t.Fatalf("Expected exactly 2 package entries in cache, got: %d", len(state.TypesInfo))
	}

	testDirAbs, err := filepath.Abs(filepath.Clean("../test"))
	if err != nil {
		t.Fatalf("Expected no error resolving abs path, got: %v", err)
	}
	helperDirAbs, err := filepath.Abs(filepath.Clean("../test/GL10_helper"))
	if err != nil {
		t.Fatalf("Expected no error resolving abs path, got: %v", err)
	}

	if state.TypesInfo[filepath.Clean(testDirAbs)] == nil {
		t.Fatal("Expected cached types info for ../test")
	}
	if state.TypesInfo[filepath.Clean(helperDirAbs)] == nil {
		t.Fatal("Expected cached types info for ../test/GL10_helper")
	}

	t.Logf("Types info cache entries: %d", len(state.TypesInfo))
}

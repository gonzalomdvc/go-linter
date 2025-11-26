package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/ast"
	"github.com/gonzalomdvc/go-linter/ui"
)

func Test_GL1(t *testing.T) {
	astFile, _, err := ast.GetAst("../subpackage/main.go")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	findings := GL1(astFile)
	if len(findings) != 1 {
		t.Errorf("Expected 1 finding, got %d", len(findings))
	}

	for _, finding := range findings {
		ui.PrintPosition(finding.Position, 20, "../subpackage")
	}
}
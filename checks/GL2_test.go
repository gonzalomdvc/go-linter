package checks

import (
	"fmt"
	"testing"

	"github.com/gonzalomdvc/go-linter/ast"
	"github.com/gonzalomdvc/go-linter/ui"
)

func Test_GL2(t *testing.T) {
	astFile, fset, err := ast.GetAst("../test/GL2.go")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	findings := GL2(fset, astFile)
	if len(findings) != 2 {
		t.Errorf("Expected 2 findings, got %d", len(findings))
	}

	for _, finding := range findings {
		pos, err := ui.PrintPosition(finding.Position, finding.Check.Message)
		if err != nil {
			t.Errorf("Error printing position: %v", err)
		} else {
			fmt.Printf("%s\n", pos)
		}
	}
}

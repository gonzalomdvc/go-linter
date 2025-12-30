package test

import (
	"fmt"
	goast "go/ast"
	"go/token"

	"github.com/gonzalomdvc/go-linter/ast"
	"github.com/gonzalomdvc/go-linter/interfaces"
	"github.com/gonzalomdvc/go-linter/ui"
)

func RunCheckTest(filename string, verbose bool, positions []interfaces.Position, checkFunc func(fset *token.FileSet, astFile *goast.File) []interfaces.Finding) error {
	astFile, fset, err := ast.GetAst(fmt.Sprintf("../test/%s", filename))
	if err != nil {
		return fmt.Errorf("Expected no error, got %v", err)
	}

	findings := checkFunc(fset, astFile)

	foundPositions := make(map[interfaces.Position]bool)
	for _, pos := range positions {
		foundPositions[pos] = false
	}
	for _, finding := range findings {
		pos := interfaces.Position{
			Column: finding.Position.Column,
			Line:   finding.Position.Line,
		}
		fmt.Printf("Detected finding at position: Column: %d, Line: %d\n", pos.Column, pos.Line)
		foundPositions[pos] = true
	}
	for pos := range foundPositions {
		if foundPositions[pos] == false {
			return fmt.Errorf("Instance of linter warning undetected at position: Column: %d, Line: %d", pos.Column, pos.Line)
		}
	}
	if verbose {
		err = ui.PrintFindings(findings)
		if err != nil {
			return fmt.Errorf("Expected no error printing findings, got %v", err)
		}
	}
	return nil
}

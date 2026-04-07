package checks

import (
	"fmt"
	"slices"

	"github.com/gonzalomdvc/go-linter/ast"
	"github.com/gonzalomdvc/go-linter/packages"
	"github.com/gonzalomdvc/go-linter/ui"
)

type Position struct {
	Column int
	Line   int
}

func RunCheckTest(filename string, verbose bool, positions []Position, checkFunc CheckFunc, state *packages.State) error {
	astFile, fset, err := ast.GetAst(fmt.Sprintf("../test/%s", filename))
	if err != nil {
		return fmt.Errorf("Expected no error, got %v", err)
	}

	findings := checkFunc(fset, astFile, state)

	foundPositions := make(map[Position]bool)
	for _, pos := range positions {
		foundPositions[pos] = false
	}
	for _, finding := range findings {
		pos := Position{
			Column: finding.Position.Column,
			Line:   finding.Position.Line,
		}
		if !slices.Contains(positions, pos) {
			return fmt.Errorf("Unexpected finding at position: Column: %d, Line: %d", pos.Column, pos.Line)
		} else {
			foundPositions[pos] = true
		}
	}
	for pos := range foundPositions {
		if foundPositions[pos] == false {
			return fmt.Errorf("Instance of linter warning undetected at position: Column: %d, Line: %d", pos.Column, pos.Line)
		}
	}
	if verbose {
		err = ui.Printfindings(findings)
		if err != nil {
			return fmt.Errorf("Expected no error printing findings, got %v", err)
		}
	}
	return nil
}

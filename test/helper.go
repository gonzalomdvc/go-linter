package test

import (
	"fmt"
	goast "go/ast"
	"go/token"

	"github.com/gonzalomdvc/go-linter/ast"
	"github.com/gonzalomdvc/go-linter/interfaces"
	"github.com/gonzalomdvc/go-linter/ui"
)

func RunCheckTest(filename string, verbose bool, checkFunc func(fset *token.FileSet, astFile *goast.File) []interfaces.Finding) ([]interfaces.Finding, error) {
	astFile, fset, err := ast.GetAst(fmt.Sprintf("../test/%s", filename))
	if err != nil {
		return nil, fmt.Errorf("Expected no error, got %v", err)
	}

	findings := checkFunc(fset, astFile)
	if verbose {
		err = ui.PrintFindings(findings)
		if err != nil {
			return nil, fmt.Errorf("Expected no error printing findings, got %v", err)
		}
	}
	return findings, nil
}

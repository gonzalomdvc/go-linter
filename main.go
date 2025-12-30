package main

import (
	"go/ast"
	"go/token"
	"os"

	"github.com/gonzalomdvc/go-linter/checks"
	"github.com/gonzalomdvc/go-linter/interfaces"
	"github.com/gonzalomdvc/go-linter/ui"
)

var Checks = []func(*token.FileSet, *ast.File) []interfaces.Finding{
	checks.GL1,
	checks.GL2,
	checks.GL3,
	checks.GL4,
	checks.GL6,
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	findings := RunLinterChecks(dir, Checks, 0)
	if len(findings) > 0 {
		ui.PrintFindings(findings)
		os.Exit(1)
	} else {
		ui.PrintSuccessfulMessage()
		os.Exit(0)
	}
}

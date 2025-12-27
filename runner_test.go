package main

import (
	"fmt"
	goast "go/ast"
	"go/token"
	"testing"

	"github.com/gonzalomdvc/go-linter/checks"
	"github.com/gonzalomdvc/go-linter/interfaces"
)

func BenchmarkRunLinterChecks(b *testing.B) {
	var checkFuncs = []func(*token.FileSet, *goast.File) []interfaces.Finding{
		checks.GL1,
		checks.GL2,
		checks.GL3,
		checks.GL4,
	}
	for cnt := 1; cnt <= len(checkFuncs); cnt++ {
		n := cnt
		b.Run(fmt.Sprintf("Running with %d checks", n), func(b *testing.B) {
			for b.Loop() {
				RunLinterChecks("./test", checkFuncs[:n], 0)
			}
		})
	}

}

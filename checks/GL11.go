// Go Linter 11: function with too many arguments
package checks

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/gonzalomdvc/go-linter/model"
	"github.com/gonzalomdvc/go-linter/packages"
)

var maxNumberofParams = 10

func GL11(fset *token.FileSet, file *ast.File, state *packages.State) []model.Finding {
	findings := []model.Finding{}
	ast.Inspect(file, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			var numParams int
			for _, l := range funcDecl.Type.Params.List {
				numParams += len(l.Names)
			}
			if numParams > maxNumberofParams {
				findings = append(findings, model.Finding{
					Position: fset.Position(funcDecl.Pos()),
					Message:  fmt.Sprintf("Function %s has too many parameters, keep under 10", funcDecl.Name),
				})
				return true
			}
		}
		return true
	})
	return findings
}

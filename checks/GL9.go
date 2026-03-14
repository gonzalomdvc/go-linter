// Go Linter X: explanation

package checks

import (
	"go/ast"
	"go/token"

	"github.com/gonzalomdvc/go-linter/model"
	"github.com/gonzalomdvc/go-linter/packages"
)

func GL9(fset *token.FileSet, file *ast.File, state *packages.State) []model.Finding {
	var findings []model.Finding
	ast.Inspect(file, func(n ast.Node) bool {
		if selectStmt, ok := n.(*ast.SelectStmt); ok {
			if len(selectStmt.Body.List) == 1 {
				findings = append(findings, model.Finding{
					Position: fset.Position(selectStmt.Pos()),
					Message:  "Should use channel receive instead of single case select.",
				})
			}
		}
		return true
	})
	return findings
}

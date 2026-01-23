// Go Linter X: explanation

package checks

import (
	"go/ast"
	"go/token"

	"github.com/gonzalomdvc/go-linter/interfaces"
)

func GL9(fset *token.FileSet, file *ast.File) []interfaces.Finding {
	var findings []interfaces.Finding
	ast.Inspect(file, func(n ast.Node) bool {
		if selectStmt, ok := n.(*ast.SelectStmt); ok {
			if len(selectStmt.Body.List) == 1 {
				findings = append(findings, interfaces.Finding{
					Position: fset.Position(selectStmt.Pos()),
					Check: interfaces.Check{
						Name:    "GL9",
						Func:    GL9,
						Message: "Should use channel receive instead of single case select.",
					},
				})
			}
		}
		return true
	})
	return findings
}

// Go Linter X: explanation

package checks

import (
	"go/ast"
	"go/token"

	"github.com/gonzalomdvc/go-linter/interfaces"
)

func GL8(fset *token.FileSet, file *ast.File) []interfaces.Finding {
	var findings []interfaces.Finding
	ast.Inspect(file, func(n ast.Node) bool {
		if _, ok := n.(*ast.ForStmt); ok {
			ast.Inspect(n, func(n ast.Node) bool {
				if selectStmt, ok := n.(*ast.SelectStmt); ok {
					for _, clause := range selectStmt.Body.List {
						if commClause, ok := clause.(*ast.CommClause); ok {
							if commClause.Comm == nil && commClause.Body == nil {
								findings = append(findings, interfaces.Finding{
									Position: fset.Position(commClause.Pos()),
									Check: interfaces.Check{
										Name:    "GL8",
										Func:    GL8,
										Message: "Empty default in for-select spins (bad for your CPU)",
									},
								})
								return true
							}
						}
					}
				}
				return true
			})
		}
		return true
	})
	return findings
}

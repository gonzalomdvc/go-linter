// Go Linter 12: don't use naked returns

package checks

import (
	"go/ast"
	"go/token"

	"github.com/gonzalomdvc/go-linter/model"
	"github.com/gonzalomdvc/go-linter/packages"
)

func GL12(fset *token.FileSet, file *ast.File, state *packages.State) []model.Finding {
	var findings []model.Finding
	ast.Inspect(file, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			hasNamedReturns := false
			if funcDecl.Type.Results == nil {
				return true
			}
			for _, param := range funcDecl.Type.Results.List {
				if param.Names != nil {
					hasNamedReturns = true
				}
			}
			if funcDecl.Body == nil {
				return true
			}
			ast.Inspect(funcDecl.Body, func(m ast.Node) bool {
				if returnStmt, ok := m.(*ast.ReturnStmt); ok {
					if returnStmt.Results == nil && hasNamedReturns {
						findings = append(findings, model.Finding{
							Position: fset.Position(returnStmt.Pos()),
							Message:  "Don't use naked returns for return parameters, as it hurts readability",
						})
						return false
					}
					return true
				}
				return true
			})

			return true
		}

		return true
	})

	return findings
}

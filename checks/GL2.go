// Go Linter 2: empty body in an if or else branch

package checks

import (
	"go/ast"
	"go/token"

	"github.com/gonzalomdvc/go-linter/model"
	"github.com/gonzalomdvc/go-linter/packages"
)

func GL2(fset *token.FileSet, file *ast.File, state *packages.State) []model.Finding {
	var findings []model.Finding

	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.IfStmt:
			if len(x.Body.List) == 0 {
				findings = append(findings, model.Finding{
					Position: fset.Position(n.Pos()),
					Message:  "Empty if statement body",
				})
			}

			if x.Else != nil {
				if elseBlock, ok := x.Else.(*ast.BlockStmt); ok {
					if len(elseBlock.List) == 0 {
						findings = append(findings, model.Finding{
							Position: fset.Position(x.Else.Pos()),
							Message:  "Empty else statement body",
						})

					}
				}
			}
		}
		return true
	})
	return findings
}

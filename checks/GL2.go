// Go Linter 2: empty body in an if or else branch

package checks

import (
	"go/ast"
	"go/token"

	"github.com/gonzalomdvc/go-linter/interfaces"
)

func GL2(fset *token.FileSet, file *ast.File) []interfaces.Finding {
	var findings []interfaces.Finding

	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.IfStmt:
			if len(x.Body.List) == 0 {
				findings = append(findings, interfaces.Finding{
					Position: fset.Position(n.Pos()), // Simplified for example
					Check: interfaces.Check{
						Name:    "GL2",
						Func:    GL2,
						Message: "Empty if statement body",
					},
				})
			}

			if x.Else != nil {
				if elseBlock, ok := x.Else.(*ast.BlockStmt); ok {
					if len(elseBlock.List) == 0 {
						findings = append(findings, interfaces.Finding{
							Position: fset.Position(x.Else.Pos()), // Simplified for example
							Check: interfaces.Check{
								Name:    "GL2",
								Func:    GL2,
								Message: "Empty else statement body",
							},
						})
					}
				}
			}
		}
		return true
	})
	return findings
}

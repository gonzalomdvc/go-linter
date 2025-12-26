// Go Linter 4: comparison of same literal always true

package checks

import (
	"go/ast"
	"go/token"

	"github.com/gonzalomdvc/go-linter/interfaces"
)

func GL4(fset *token.FileSet, file *ast.File) []interfaces.Finding {
	var findings []interfaces.Finding
	ast.Inspect(file, func(n ast.Node) bool {
		bop, ok := n.(*ast.BinaryExpr)
		if !ok || bop.Op != token.EQL {
			return true
		}
		xlit, ok := bop.X.(*ast.BasicLit)
		if !ok {
			return true
		}
		ylit, ok := bop.Y.(*ast.BasicLit)
		if !ok {
			return true
		}
		if xlit.Value == ylit.Value {
			findings = append(findings, interfaces.Finding{
				Position: fset.Position(bop.Pos()),
				Check: interfaces.Check{
					Name:    "GL4",
					Func:    GL4,
					Message: "Same expression on both sides of operand ==",
				},
			})
		}
		return true
	})

	return findings
}

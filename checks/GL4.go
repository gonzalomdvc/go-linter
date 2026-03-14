// Go Linter 4: comparison of same literal always true

package checks

import (
	"go/ast"
	"go/token"

	"github.com/gonzalomdvc/go-linter/model"
	"github.com/gonzalomdvc/go-linter/packages"
)

func GL4(fset *token.FileSet, file *ast.File, state *packages.State) []model.Finding {
	var findings []model.Finding
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
			findings = append(findings, model.Finding{
				Position: fset.Position(bop.Pos()),
				Message:  "Same expression on both sides of operand ==",
			})
		}
		return true
	})

	return findings
}

// Go Linter 6: if else should be written as cases

package checks

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/gonzalomdvc/go-linter/interfaces"
)

func GL6(fset *token.FileSet, file *ast.File) []interfaces.Finding {
	var findings []interfaces.Finding
	var checkedIfs map[token.Pos]bool = make(map[token.Pos]bool)
	ast.Inspect(file, func(n ast.Node) bool {
		var variable *ast.Ident
		if ifs, ok := n.(*ast.IfStmt); ok {
			if ifs.Else == nil {
				return false
			}
			// We have to keep track of checked ifs to ignore them when traversing the AST normally
			if _, alreadyChecked := checkedIfs[ifs.Pos()]; alreadyChecked {
				return false
			}
			for {
				if bin, ok := ifs.Cond.(*ast.BinaryExpr); ok {
					if xvar, ok := bin.X.(*ast.Ident); ok {
						if variable == nil {
							variable = xvar
						} else {
							if variable.Name != xvar.Name {
								return false
							}
						}
					}
				}
				elseIf, ok := ifs.Else.(*ast.IfStmt)
				if !ok {
					break
				}
				checkedIfs[ifs.Pos()] = true
				ifs = elseIf
			}

			if variable != nil {
				findings = append(findings, interfaces.Finding{
					Position: fset.Position(n.Pos()),
					Check: interfaces.Check{
						Name:    "GL6",
						Func:    GL6,
						Message: fmt.Sprintf("If-Else statements over a single variable %s should be written as switch-case", variable.Name),
					},
				})
			}

		}
		return true
	})

	return findings
}

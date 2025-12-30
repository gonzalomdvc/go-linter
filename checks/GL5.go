// Go Linter 5: Redundant new line on Println

// Nombre de función Println y regex match \n

package checks

import (
	"go/ast"
	"go/token"

	"github.com/gonzalomdvc/go-linter/interfaces"
)

func GL5(fset *token.FileSet, file *ast.File) []interfaces.Finding {
	var findings []interfaces.Finding
	ast.Inspect(file, func(n ast.Node) bool {
		var call bool = false
		if expr, ok := n.(*ast.CallExpr); ok {
			if sel, ok := expr.Fun.(*ast.SelectorExpr); ok {
				if identX, ok1 := sel.X.(*ast.Ident); ok1 {
					if identX.Name == "fmt" && sel.Sel.Name == "Println" {
						call = true
					}
				}
			}
			if firstArg, ok := expr.Args[0].(*ast.BasicLit); ok && call {
				if firstArg.Kind == token.STRING && len(firstArg.Value) >= 2 {
					match := firstArg.Value[len(firstArg.Value)-3:len(firstArg.Value)-1] == "\\n"
					if match {
						findings = append(findings, interfaces.Finding{
							Position: fset.Position(n.Pos()),
							Check: interfaces.Check{
								Name:    "GL5",
								Func:    GL5,
								Message: "fmt.Println already adds a new line, no need to include \\n in the string",
							},
						})
					}
				}
				return true
			}
		}
		return true
	})

	return findings
}

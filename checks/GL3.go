// Go Linter 3: Infinite recursive call

package checks

import (
	"go/ast"
	"go/token"

	"github.com/gonzalomdvc/go-linter/interfaces"
)

func GL3(fset *token.FileSet, file *ast.File) []interfaces.Finding {
	var findings []interfaces.Finding
	ast.Inspect(file, func(n ast.Node) bool {
		fd, ok := n.(*ast.FuncDecl)
		if !ok || fd.Body == nil {
			return true
		}

		var selfCall *ast.CallExpr
		var hasReturn bool

		ast.Inspect(fd.Body, func(n2 ast.Node) bool {
			switch s := n2.(type) {
			case *ast.ReturnStmt:
				hasReturn = true
			case *ast.CallExpr:
				if ident, ok := s.Fun.(*ast.Ident); ok {
					if ident.Name == fd.Name.Name {
						selfCall = s
					}
				}
			}
			if hasReturn && selfCall != nil {
				return false
			}
			return true
		})

		if selfCall != nil && !hasReturn {
			findings = append(findings, interfaces.Finding{
				Position: fset.Position(selfCall.Pos()),
				Check: interfaces.Check{
					Name:    "GL3",
					Func:    GL3,
					Message: "Recursive function without exit condition",
				},
			})
		}

		return true
	})

	return findings
}

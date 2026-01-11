// Go linter 7: use camelCase in function declarations!

package checks

import (
	"go/ast"
	"go/token"
	"regexp"

	"github.com/gonzalomdvc/go-linter/interfaces"
)

func GL7(fset *token.FileSet, file *ast.File) []interfaces.Finding {
	var findings []interfaces.Finding

	ast.Inspect(file, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			matched, err := regexp.Match("^[a-z]+_[a-z]+$", []byte(funcDecl.Name.Name))
			if err == nil && matched {
				findings = append(findings, interfaces.Finding{
					Position: fset.Position(n.Pos()),
					Check: interfaces.Check{
						Name:    "GL7",
						Func:    GL7,
						Message: "Function names should use camelCase instead of snake_case",
					},
				})
			}
		}
		return true
	})
	return findings
}

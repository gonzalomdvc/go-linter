// Go linter 7: use camelCase in function declarations!

package checks

import (
	"go/ast"
	"go/token"
	"regexp"

	"github.com/gonzalomdvc/go-linter/model"
	"github.com/gonzalomdvc/go-linter/packages"
)

func GL7(fset *token.FileSet, file *ast.File, state *packages.State) []model.Finding {
	var findings []model.Finding

	ast.Inspect(file, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			matched, err := regexp.Match("^[a-z]+_[a-z]+$", []byte(funcDecl.Name.Name))
			if err == nil && matched {
				findings = append(findings, model.Finding{
					Position: fset.Position(n.Pos()),
					Message:  "Function names should use camelCase instead of snake_case",
				})
			}
		}
		return true
	})
	return findings
}

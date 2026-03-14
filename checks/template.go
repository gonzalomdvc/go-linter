// Go Linter X: explanation

package checks

import (
	"go/ast"
	"go/token"

	"github.com/gonzalomdvc/go-linter/model"
	"github.com/gonzalomdvc/go-linter/packages"
)

func GLX(fset *token.FileSet, file *ast.File, state *packages.State) []model.Finding {
	ast.Inspect(file, func(n ast.Node) bool {
		// Implementation
		return true
	})
	// Implementation
	return []model.Finding{}
}

// Go Linter X: explanation

package checks

import (
	"go/ast"
	"go/token"

	"github.com/gonzalomdvc/go-linter/interfaces"
)

func GLX(fset *token.FileSet, file *ast.File) []interfaces.Finding {
	ast.Inspect(file, func(n ast.Node) bool {
		// Implementation
		return true
	})
	// Implementation
	return []interfaces.Finding{}
}

package checks

import (
	"go/ast"
	"go/token"

	"github.com/gonzalomdvc/go-linter/model"
	"github.com/gonzalomdvc/go-linter/packages"
)

type Check struct {
	Name    string
	Func    CheckFunc
	Message string
}

type CheckFunc func(fset *token.FileSet, astFile *ast.File, state *packages.State) []model.Finding

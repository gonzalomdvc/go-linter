package interfaces

import (
	"go/ast"
	"go/token"
)

type Check struct {
	Name    string
	Func    func(*token.FileSet, *ast.File) []Finding
	Message string
}

type Finding struct {
	Position token.Position
	Check    Check
}

type Position struct {
	Column int
	Line   int
}

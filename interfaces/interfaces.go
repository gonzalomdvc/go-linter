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

// Position represents a specific location in a source file. Offsets in chars.
// Note: this might not match the column and line you see on your text editor.
// This is because tabs are counted as a single character here.
type Position struct {
	Column int
	Line   int
}

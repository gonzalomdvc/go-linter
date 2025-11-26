package interfaces

import (
	"go/ast"
	"go/token"
)

type Check struct {
	Name string
	Func func(*ast.File) ([]Finding)
	Message string
} 

type Finding struct {
	Position token.Position
	Check Check
}
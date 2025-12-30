package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func GetAst(fileName string) (*ast.File, *token.FileSet, error) {
	srcBytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, nil, err
	}
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, fileName, srcBytes, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}
	return astFile, fset, nil
}

func PrintAst(fset *token.FileSet, astFile *ast.File) {
	ast.Print(fset, astFile)
}

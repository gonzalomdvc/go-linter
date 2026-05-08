package ast

import "testing"

func Test_PrintAst(t *testing.T) {
	astFile, fset, _ := GetAst("../test/GL14.go", nil)
	PrintAst(fset, astFile)
}

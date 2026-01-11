package ast

import "testing"

func Test_PrintAst(t *testing.T) {
	astFile, fset, _ := GetAst("../main.go")
	PrintAst(fset, astFile)
}

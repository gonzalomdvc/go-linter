// Go Linter 14: Consider using 'rune' or 'byte' instead of 'int' for character data

package checks

import (
	"go/ast"
	"go/token"

	"github.com/gonzalomdvc/go-linter/model"
	"github.com/gonzalomdvc/go-linter/packages"
)

func GL14(fset *token.FileSet, file *ast.File, state *packages.State) []model.Finding {
	var findings []model.Finding

	ast.Inspect(file, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					vType, ok := valueSpec.Type.(*ast.Ident)
					if !ok {
						return true
					}
					if isRuneOrCharUnderlyingType(vType) && hasCharValue(valueSpec.Values) {
						findings = append(findings, model.Finding{
							Position: fset.Position(genDecl.TokPos),
							Message:  "Consider using alias 'rune' for int32 chars and 'byte' for uint8 chars",
						})
					}
				}
			}
		}
		return true
	})
	return findings
}

func isRuneOrCharUnderlyingType(vType *ast.Ident) bool {
	if vType.Name == "int32" || vType.Name == "uint8" {
		return true
	} else {
		return false
	}
}

func hasCharValue(values []ast.Expr) bool {
	for _, expr := range values {
		if basicLit, ok := expr.(*ast.BasicLit); ok {
			if basicLit.Kind == token.CHAR {
				return true
			}
		}
	}

	return false
}

// Go Linter 14: Consider using 'rune' or 'byte' instead of 'int' for character data

package checks

import (
	"go/ast"
	"go/token"
	gotypes "go/types"

	"github.com/gonzalomdvc/go-linter/model"
	"github.com/gonzalomdvc/go-linter/packages"
	"github.com/gonzalomdvc/go-linter/types"
)

func GL14(fset *token.FileSet, file *ast.File, state *packages.State) []model.Finding {
	var findings []model.Finding

	// Hook for future type-driven checks. Current lint behavior remains unchanged.

	info, pkg := types.GetTypesInfoForFile(fset, file, state)

	// Build a position→Object map once to avoid O(n) scan per ident.
	posToObj := make(map[token.Pos]gotypes.Object, len(info.Uses))
	for id, obj := range info.Uses {
		posToObj[id.Pos()] = obj
	}

	ast.Inspect(file, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					vType, ok := valueSpec.Type.(*ast.Ident)
					if !ok {
						return true
					}
					if isRuneOrCharUnderlyingType(vType, pkg, info) && hasCharValue(valueSpec.Values) {
						findings = append(findings, model.Finding{
							Position: fset.Position(genDecl.TokPos),
							Message:  "Consider using alias 'rune' for int32 chars and 'byte' for uint8 chars",
						})
					}
				}
			}
		}
		if assignStmt, ok := n.(*ast.AssignStmt); ok {
			for _, lhs := range assignStmt.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok {
					if isAssignRuneOrCharType(ident, posToObj) && hasCharValue(assignStmt.Rhs) {
						findings = append(findings, model.Finding{
							Position: fset.Position(assignStmt.Pos()),
							Message:  "Consider using alias 'rune' for int32 chars and 'byte' for uint8 chars",
						})
					}
					return true
				}
			}
		}

		return true
	})
	return findings
}

func isRuneOrCharUnderlyingType(vType *ast.Ident, pkg *gotypes.Package, info *gotypes.Info) bool {
	if vType.Name == "int32" || vType.Name == "uint8" {
		return true
	}

	if obj := pkg.Scope().Lookup(vType.Name); obj != nil {
		if obj.Type().String() == "int32" || obj.Type().String() == "uint8" {
			return true
		}
	}

	return false
}

func isAssignRuneOrCharType(ident *ast.Ident, posToObj map[token.Pos]gotypes.Object) bool {
	if obj, ok := posToObj[ident.Pos()]; ok {
		t := obj.Type().String()
		return t == "int32" || t == "uint8"
	}
	return false
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

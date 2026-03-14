// Go Linter 1: no unused local vars

package checks

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/gonzalomdvc/go-linter/model"
	"github.com/gonzalomdvc/go-linter/packages"
)

func GL1(fset *token.FileSet, file *ast.File, state *packages.State) []model.Finding {
	var findings []model.Finding
	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.DeclStmt:
			if genDecl, ok := x.Decl.(*ast.GenDecl); ok {
				for _, spec := range genDecl.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						for _, name := range valueSpec.Names {
							if name.Obj != nil && name.Obj.Kind == ast.Var {
								used := false
								ast.Inspect(file, func(n2 ast.Node) bool {
									if ident, ok := n2.(*ast.Ident); ok {
										if ident.Name == name.Name && ident != name {
											used = true
											return false
										}
									}
									return true
								})
								if !used {
									findings = append(findings, model.Finding{
										Position: fset.Position(n.Pos()),
										Message:  fmt.Sprintf("Unused variable: %s", name.Name),
									})

								}
							}
						}
					}
				}
			}
		}

		return true
	})

	return findings
}

// Go Linter 1: no unused local vars

package checks

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/gonzalomdvc/go-linter/interfaces"
)

func GL1(fset *token.FileSet, file *ast.File) []interfaces.Finding {
	var findings []interfaces.Finding
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
											findings = append(findings, interfaces.Finding{
												Position: fset.Position(n.Pos()), // Simplified for example
												Check: interfaces.Check{
													Name: "GL1",
													Func: GL1,
													Message: fmt.Sprintf("Unused variable: %s", name.Name),
												},
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
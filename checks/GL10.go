// Go Linter 10: using a Deprecated function

package checks

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"strings"

	"github.com/gonzalomdvc/go-linter/interfaces"
	"golang.org/x/tools/go/packages"
)

func GL10(fset *token.FileSet, file *ast.File) []interfaces.Finding {
	findings := []interfaces.Finding{}
	packageAddresses := []string{}
	callExprs := map[string]*ast.CallExpr{}
	ast.Inspect(file, func(n ast.Node) bool {
		im, ok := n.(*ast.ImportSpec)
		if ok {
			packageAddresses = append(packageAddresses, strings.Trim(im.Path.Value, `"`))
		}
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		callExprs[selExpr.Sel.Name] = callExpr
		return true
	})
	cfg := &packages.Config{Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo}

	pkgs, err := packages.Load(cfg, packageAddresses...)
	if err != nil {
		log.Fatal(err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		log.Fatal("package load errors")
	}
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				decl, ok := n.(*ast.FuncDecl)
				if !ok {
					return true
				}
				if decl.Doc == nil {
					return true
				}
				matched := callExprs[decl.Name.Name]
				if matched == nil {
					return true
				}
				for _, comment := range decl.Doc.List {
					if strings.Contains(comment.Text, "Deprecated") {
						findings = append(findings, interfaces.Finding{
							Position: fset.Position(matched.Pos()),
							Check: interfaces.Check{
								Name:    "GL10",
								Func:    GL10,
								Message: fmt.Sprintf("Function %s is deprecated: %s", decl.Name.Name, strings.TrimSpace(strings.Split(comment.Text, "Deprecated:")[1])),
							},
						})
						break
					}
				}
				return true
			})
		}
	}

	return findings
}

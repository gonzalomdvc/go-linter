// Go Linter 10: using a Deprecated function

package checks

import (
	"fmt"
	"go/ast"
	"go/token"
	"slices"
	"strings"

	"github.com/gonzalomdvc/go-linter/model"
	"github.com/gonzalomdvc/go-linter/packages"
)

func GL10(fset *token.FileSet, file *ast.File, state *packages.State) []model.Finding {
	findings := []model.Finding{}
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

	pkgs := state.Packages
	for key, pkg := range pkgs {
		if !slices.Contains(packageAddresses, key) {
			continue
		}
		for _, decl := range pkg.FuncDecls {
			if decl.Doc == nil {
				continue
			}
			matched := callExprs[decl.Name.Name]
			if matched == nil {
				continue
			}
			for _, comment := range decl.Doc.List {
				if strings.Contains(comment.Text, "Deprecated") {
					findings = append(findings, model.Finding{
						Position: fset.Position(matched.Pos()),
						Message:  fmt.Sprintf("Function %s is deprecated: %s", decl.Name.Name, strings.TrimSpace(strings.Split(comment.Text, "Deprecated:")[1])),
					})
					break
				}

			}
		}
	}

	return findings
}

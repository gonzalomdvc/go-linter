package checks

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"github.com/gonzalomdvc/go-linter/model"
	"github.com/gonzalomdvc/go-linter/packages"
	gotypes "github.com/gonzalomdvc/go-linter/types"
)

func GL10(fset *token.FileSet, file *ast.File, state *packages.State) []model.Finding {
	findings := []model.Finding{}

	info, _ := gotypes.GetTypesInfoForFile(fset, file, state)
	if info == nil {
		// Fallback: no types info available (tests may not populate it).
		// Build a map of import alias -> import path from the AST, then
		// match selector calls against the imported package declarations
		// stored in state.Packages.
		aliasToPath := make(map[string]string)
		for _, im := range file.Imports {
			trimmedPath := strings.Trim(im.Path.Value, `"`)
			name := ""
			if im.Name != nil {
				name = im.Name.Name
			} else {
				parts := strings.Split(trimmedPath, "/")
				name = parts[len(parts)-1]
			}
			aliasToPath[name] = trimmedPath
		}

		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			xIdent, ok := sel.X.(*ast.Ident)
			if !ok {
				return true
			}

			importPath, ok := aliasToPath[xIdent.Name]
			if !ok {
				return true
			}

			pkg, ok := state.Packages[importPath]
			if !ok {
				return true
			}

			for _, decl := range pkg.FuncDecls {
				if decl.Name.Name != sel.Sel.Name || decl.Doc == nil {
					continue
				}

				for _, comment := range decl.Doc.List {
					if !strings.Contains(comment.Text, "Deprecated:") {
						continue
					}

					msg := strings.TrimSpace(strings.SplitN(comment.Text, "Deprecated:", 2)[1])

					findings = append(findings, model.Finding{
						Position: fset.Position(sel.Pos()),
						Message:  fmt.Sprintf("Function %s is deprecated: %s", sel.Sel.Name, msg),
					})

					return true
				}

				break
			}

			return true
		})
		return findings
	}

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		obj := info.Uses[sel.Sel]

		fn, ok := obj.(*types.Func)
		if !ok || fn.Pkg() == nil {
			return true
		}

		pkg, ok := state.Packages[fn.Pkg().Path()]
		if !ok {
			return true
		}

		for _, decl := range pkg.FuncDecls {
			if decl.Name.Name != fn.Name() || decl.Doc == nil {
				continue
			}

			for _, comment := range decl.Doc.List {
				if !strings.Contains(comment.Text, "Deprecated:") {
					continue
				}

				msg := strings.TrimSpace(strings.SplitN(comment.Text, "Deprecated:", 2)[1])

				findings = append(findings, model.Finding{
					Position: fset.Position(sel.Pos()),
					Message:  fmt.Sprintf("Function %s is deprecated: %s", fn.Name(), msg),
				})

				return true
			}

			break
		}

		return true
	})

	return findings
}

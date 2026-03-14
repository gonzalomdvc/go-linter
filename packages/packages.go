package packages

import (
	"fmt"
	goast "go/ast"
	"go/token"
	"strings"

	"github.com/gonzalomdvc/go-linter/ast"
	"golang.org/x/tools/go/packages"
)

type FuncDeclResult struct {
	PackagePath string
	FuncDecls   []*goast.FuncDecl
}

type Package struct {
	FuncDecls []*goast.FuncDecl
}

type SourceAst struct {
	Fset    *token.FileSet
	AstFile *goast.File
}

type State struct {
	SourceAsts map[string]SourceAst
	Packages   map[string]Package
}

func ImportPackages(astFile *goast.File, funcDecls chan FuncDeclResult, state *State) {
	packageAddresses := []string{}
	cfg := &packages.Config{Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax}
	goast.Inspect(astFile, func(n goast.Node) bool {
		im, ok := n.(*goast.ImportSpec)
		if ok {
			trimmedPath := strings.Trim(im.Path.Value, `"`)
			if state.Packages[trimmedPath].FuncDecls != nil {
				return true
			}
			packageAddresses = append(packageAddresses, trimmedPath)
		}

		return true
	})

	pkgs, err := packages.Load(cfg, packageAddresses...)
	if err != nil {
		panic(fmt.Sprintf("Error loading packages: %s", err))
	}

	for _, pkg := range pkgs {
		funcDecls <- FuncDeclResult{PackagePath: pkg.PkgPath, FuncDecls: ast.GetFuncDecls(pkg.Syntax)}
	}

}

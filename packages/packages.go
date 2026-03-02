package packages

import (
	"fmt"
	goast "go/ast"
	"strings"

	"github.com/gonzalomdvc/go-linter/ast"
	"golang.org/x/tools/go/packages"
)

type FuncDeclResult struct {
	PackagePath string
	FuncDecls   []*goast.FuncDecl
}

func ImportPackages(astFile *goast.File, funcDecls chan FuncDeclResult) {
	packageAddresses := []string{}
	cfg := &packages.Config{Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo}
	goast.Inspect(astFile, func(n goast.Node) bool {
		im, ok := n.(*goast.ImportSpec)
		if ok {
			packageAddresses = append(packageAddresses, strings.Trim(im.Path.Value, `"`))
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

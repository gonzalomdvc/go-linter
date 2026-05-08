package packages

import (
	"fmt"
	goast "go/ast"
	"go/token"
	"os"
	"strings"

	"go/types"

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
	SourceAsts       map[string]SourceAst
	Packages         map[string]Package
	TypesInfo        map[string]*types.Info
	PackageTypesInfo map[string]*packages.Package
}

func collectImportPaths(astFile *goast.File, state *State, seen map[string]struct{}, importPaths *[]string) {
	goast.Inspect(astFile, func(n goast.Node) bool {
		im, ok := n.(*goast.ImportSpec)
		if !ok {
			return true
		}

		trimmedPath := strings.Trim(im.Path.Value, `"`)
		if state.Packages[trimmedPath].FuncDecls != nil {
			return true
		}
		if _, exists := seen[trimmedPath]; exists {
			return true
		}

		seen[trimmedPath] = struct{}{}
		*importPaths = append(*importPaths, trimmedPath)
		return true
	})
}

func getPackageLoadConfig() *packages.Config {
	return &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedModule,
		Env: append(os.Environ(),
			"GOPROXY=off", // disable downloading modules
			"GOSUMDB=off", // avoid checksum DB network calls
		),
		BuildFlags: []string{"-mod=readonly"}, // do not edit go.mod/go.sum
	}
}

func loadPackagesByImportPath(importPaths []string) ([]FuncDeclResult, error) {
	if len(importPaths) == 0 {
		return nil, nil
	}

	pkgs, err := packages.Load(getPackageLoadConfig(), importPaths...)
	if err != nil {
		return nil, err
	}

	results := make([]FuncDeclResult, 0, len(pkgs))
	for _, pkg := range pkgs {
		results = append(results, FuncDeclResult{PackagePath: pkg.PkgPath, FuncDecls: ast.GetFuncDecls(pkg.Syntax)})
	}

	return results, nil
}

func ImportPackagesFromState(state *State) ([]FuncDeclResult, error) {
	seen := make(map[string]struct{})
	importPaths := make([]string, 0)

	for _, sourceAst := range state.SourceAsts {
		collectImportPaths(sourceAst.AstFile, state, seen, &importPaths)
	}

	return loadPackagesByImportPath(importPaths)
}

func ImportPackages(astFile *goast.File, funcDecls chan FuncDeclResult, state *State) {
	seen := make(map[string]struct{})
	importPaths := make([]string, 0)
	collectImportPaths(astFile, state, seen, &importPaths)

	results, err := loadPackagesByImportPath(importPaths)
	if err != nil {
		fmt.Print("Packages not found in mod cache. Run go mod tidy to download")
		return
	}

	for _, result := range results {
		funcDecls <- result
	}

}

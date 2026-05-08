// Package types interacts with go/types to extract type information from Go source code.
package types

import (
	"errors"
	"fmt"
	"go/ast"
	"go/importer"
	"go/token"
	"go/types"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gonzalomdvc/go-linter/packages"
	xpackages "golang.org/x/tools/go/packages"
)

func normalizeDirPath(dirPath string) string {
	cleaned := filepath.Clean(dirPath)
	absPath, err := filepath.Abs(cleaned)
	if err != nil {
		return cleaned
	}
	return filepath.Clean(absPath)
}

func resolveTypesFromAst(path string, src packages.SourceAst, state *packages.State, fset *token.FileSet) (*types.Info, *xpackages.Package, error) {
	info, typesPkg, err := resolveWithTypesCheck(path, state)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to resolve types for directory %s: %w", path, err)
	}

	pkg := &xpackages.Package{Types: typesPkg}
	return info, pkg, nil

}

func resolveWithTypesCheck(dirPath string, state *packages.State) (*types.Info, *types.Package, error) {
	var astFiles []*ast.File
	var sharedFset *token.FileSet
	for filePath, src := range state.SourceAsts {
		if normalizeDirPath(filepath.Dir(filePath)) == dirPath {
			astFiles = append(astFiles, src.AstFile)
			sharedFset = src.Fset
		}
	}
	if len(astFiles) == 0 {
		return nil, nil, fmt.Errorf("no AST files found for directory %s", dirPath)
	}

	cfg := &types.Config{
		Importer: importer.Default(),
		Error:    func(err error) {}, // suppress import errors; local type info is still populated
	}
	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Scopes:     make(map[ast.Node]*types.Scope),
		InitOrder:  []*types.Initializer{},
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}
	pkg, _ := cfg.Check(dirPath, sharedFset, astFiles, info)
	if pkg == nil {
		return nil, nil, fmt.Errorf("type-checking returned nil package for directory %s", dirPath)
	}
	return info, pkg, nil

}

func GetTypesInfoForFile(fset *token.FileSet, file *ast.File, state *packages.State) (*types.Info, *types.Package) {
	if fset == nil || file == nil || state == nil || state.TypesInfo == nil {
		return nil, nil
	}

	filePath := fset.Position(file.Pos()).Filename
	if filePath == "" {
		return nil, nil
	}

	dirPath := filepath.Clean(filepath.Dir(filePath))
	absDirPath, err := filepath.Abs(dirPath)
	if err == nil {
		dirPath = filepath.Clean(absDirPath)
	}
	if info, exists := state.TypesInfo[dirPath]; exists {
		return info, state.PackageTypesInfo[dirPath].Types
	}

	info, pkg, err := resolveTypesFromAst(dirPath, packages.SourceAst{Fset: fset, AstFile: file}, state, fset)
	if err != nil {
		fmt.Printf("Error resolving types for file %s: %s\n", filePath, err)
		return nil, nil
	}

	state.TypesInfo[dirPath] = info
	state.PackageTypesInfo[dirPath] = pkg
	if pkg == nil {
		return info, nil
	}
	return info, pkg.Types
}

func ResolveTypesInfoForFiles(state *packages.State, fset *token.FileSet) error {
	if state == nil {
		return errors.New("state cannot be nil")
	}
	if state.TypesInfo == nil {
		state.TypesInfo = make(map[string]*types.Info)
	}
	if state.PackageTypesInfo == nil {
		state.PackageTypesInfo = make(map[string]*xpackages.Package)
	}

	seenDirs := make(map[string]struct{})
	var errs []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	for path, ast := range state.SourceAsts {
		dir := normalizeDirPath(filepath.Dir(path))

		if _, seen := seenDirs[dir]; seen {
			continue
		}
		seenDirs[dir] = struct{}{}

		if _, exists := state.TypesInfo[dir]; exists {
			continue
		}

		wg.Add(1)
		go func(dir string, ast packages.SourceAst) {
			defer wg.Done()
			info, pkg, err := resolveTypesFromAst(dir, ast, state, fset)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errs = append(errs, fmt.Sprintf("%s: %s", path, err))
				return
			}
			state.TypesInfo[dir] = info
			state.PackageTypesInfo[dir] = pkg
		}(dir, ast)
	}

	wg.Wait()

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

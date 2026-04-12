package runner

import (
	"fmt"
	goast "go/ast"
	"go/token"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/gonzalomdvc/go-linter/ast"
	"github.com/gonzalomdvc/go-linter/checks"
	"github.com/gonzalomdvc/go-linter/model"
	"github.com/gonzalomdvc/go-linter/packages"
)

var MaxDepth = 20

var Checks = []checks.CheckFunc{
	checks.GL1,
	checks.GL2,
	checks.GL3,
	checks.GL4,
	checks.GL5,
	checks.GL6,
	checks.GL7,
	checks.GL8,
	checks.GL9,
	checks.GL10,
	checks.GL11,
	checks.GL12,
}

var ChecksNeedState = []checks.CheckFunc{
	checks.GL10,
}

func RunLinterChecks(dirname string, checkFuncs []checks.CheckFunc, depth int, parallel bool) []model.Finding {
	// Get all source files in the directory and subdirectories up to the specified depth
	srcFiles, err := getSourceFiles(dirname, depth, 0)
	if err != nil {
		panic(fmt.Sprintf("Error reading source code files: %s", err))
	}
	var findings []model.Finding

	var checksNeedingState []checks.CheckFunc
	var checksNotNeedingState []checks.CheckFunc
	for _, check := range checkFuncs {
		if contains(ChecksNeedState, check) {
			checksNeedingState = append(checksNeedingState, check)
		} else {
			checksNotNeedingState = append(checksNotNeedingState, check)
			// We can run checks that don't need state in parallel without waiting for the state to be populated
			if parallel {
				findings = append(findings, runChecksInParallel(srcFiles, []checks.CheckFunc{check}, &packages.State{Packages: make(map[string]packages.Package), SourceAsts: make(map[string]packages.SourceAst)})...)
			} else {
				findings = append(findings, runChecksSerially(srcFiles, []checks.CheckFunc{check}, &packages.State{Packages: make(map[string]packages.Package), SourceAsts: make(map[string]packages.SourceAst)})...)
			}
		}
	}

	if len(checksNeedingState) == 0 {
		return findings
	}

	// We will pass state containing auxiliary information to the checks, such as function declarations, to avoid redundant parsing and improve performance.
	var wg sync.WaitGroup
	wg.Add(len(srcFiles))
	state := &packages.State{Packages: make(map[string]packages.Package), SourceAsts: make(map[string]packages.SourceAst)}
	astFileCh := make(chan packages.SourceAst, 10)

	var consumerWg sync.WaitGroup
	consumerWg.Add(1)
	go func() {
		defer consumerWg.Done()
		for astResult := range astFileCh {
			state.SourceAsts[astResult.Fset.Position(astResult.AstFile.Pos()).Filename] = astResult
		}
	}()

	for _, filePath := range srcFiles {
		// Populate state with source files ASTs funcDecls from imported packages
		// This is for packages that need some state data depending on their AST, so that we prevent fetching the ASTs multiple times
		go func(filePath string) {
			defer wg.Done()
			astFile, fset, err := ast.GetAst(filePath)
			if err != nil {
				fmt.Printf("Error generating AST for file %s: %s\n", filePath, err)
				return
			}
			// Store the AST and FileSet in the state for later use by checks
			astFileCh <- packages.SourceAst{Fset: fset, AstFile: astFile}
		}(filePath)

	}

	// Wait for all producers, then close the channel so the consumer can finish
	wg.Wait()
	close(astFileCh)
	consumerWg.Wait()

	funcDeclResults, err := packages.ImportPackagesFromState(state)
	if err != nil {
		fmt.Print("Packages not found in mod cache. Run go mod tidy to download")
	}
	for _, funcDeclResult := range funcDeclResults {
		if _, exists := state.Packages[funcDeclResult.PackagePath]; !exists {
			state.Packages[funcDeclResult.PackagePath] = packages.Package{FuncDecls: funcDeclResult.FuncDecls}
		}
	}

	if parallel {
		findings = append(findings, runChecksInParallel(srcFiles, checkFuncs, state)...)
	} else {
		findings = append(findings, runChecksSerially(srcFiles, checkFuncs, state)...)
	}

	return findings
}

func runChecksInParallel(srcFiles []string, checkFuncs []checks.CheckFunc, state *packages.State) []model.Finding {
	var findings []model.Finding
	totalJobs := len(srcFiles) * len(checkFuncs)
	if totalJobs == 0 {
		return nil
	}

	resultsCh := make(chan []model.Finding, 10)

	for _, filePath := range srcFiles {
		go func(filePath string, state *packages.State) {
			var astFile *goast.File
			var fset *token.FileSet
			var err error
			if _, exists := state.SourceAsts[filePath]; !exists {
				astFile, fset, err = ast.GetAst(filePath)
				if err != nil {
					fmt.Printf("Error generating AST for file %s: %s\n", filePath, err)
					return
				}

			} else {
				astFile, fset = state.SourceAsts[filePath].AstFile, state.SourceAsts[filePath].Fset
			}
			for _, check := range checkFuncs {
				res := check(fset, astFile, state)
				resultsCh <- res
			}
		}(filePath, state)
	}

	for i := 0; i < totalJobs; i++ {
		res := <-resultsCh
		findings = append(findings, res...)
	}
	close(resultsCh)
	return findings
}

func runChecksSerially(srcFiles []string, checkFuncs []checks.CheckFunc, state *packages.State) []model.Finding {
	var findings []model.Finding
	for _, filePath := range srcFiles {
		var astFile *goast.File
		var fset *token.FileSet
		var err error
		if _, exists := state.SourceAsts[filePath]; !exists {
			astFile, fset, err = ast.GetAst(filePath)
			if err != nil {
				fmt.Printf("Error generating AST for file %s: %s\n", filePath, err)
				continue
			}
		} else {
			astFile, fset = state.SourceAsts[filePath].AstFile, state.SourceAsts[filePath].Fset
		}
		for _, check := range checkFuncs {
			res := check(fset, astFile, state)
			findings = append(findings, res...)
		}

	}
	return findings
}

func getSourceFiles(dirname string, depth, currentDepth int) ([]string, error) {
	var srcFiles []string
	files, err := os.ReadDir(dirname)
	if err != nil {
		return nil, fmt.Errorf("Error reading source code files: %s", err)
	}
	for _, file := range files {
		if strings.Contains(file.Name(), "helper") {
			continue
		}
		if file.IsDir() {
			if currentDepth > MaxDepth {
				fmt.Printf("Max depth of %d nested directories reached. Skipping directory: %s\n", MaxDepth, file.Name())
				continue
			}
			if currentDepth > depth {
				continue
			}
			if file.Name()[0] == '.' {
				continue
			}
			subDirPath := dirname + string(os.PathSeparator) + file.Name()
			subDirFiles, err := getSourceFiles(subDirPath, depth, currentDepth+1)
			if err != nil {
				fmt.Printf("Error getting source files from directory %s: %s\n", subDirPath, err)
				continue
			}
			srcFiles = append(srcFiles, subDirFiles...)
			continue
		}
		isSourceFile, err := regexp.MatchString(`\.go$`, file.Name())
		if err != nil {
			panic(fmt.Sprintf("Error matching file name: %s", err))
		}
		if isSourceFile {
			path := dirname + string(os.PathSeparator) + file.Name()
			srcFiles = append(srcFiles, path)
		}
	}
	return srcFiles, nil
}

func contains(checks []checks.CheckFunc, check checks.CheckFunc) bool {
	for _, c := range checks {
		if fmt.Sprintf("%p", c) == fmt.Sprintf("%p", check) {
			return true
		}
	}
	return false
}

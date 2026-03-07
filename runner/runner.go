package runner

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/gonzalomdvc/go-linter/ast"
	"github.com/gonzalomdvc/go-linter/checks"
	"github.com/gonzalomdvc/go-linter/interfaces"
	"github.com/gonzalomdvc/go-linter/packages"
)

var MaxDepth = 20

var Checks = []interfaces.CheckFunc{
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
}

func RunLinterChecks(dirname string, checks []interfaces.CheckFunc, depth int, currentDepth int, parallel bool) []interfaces.Finding {
	files, err := os.ReadDir(dirname)
	if err != nil {
		panic(fmt.Sprintf("Error reading source code files: %s", err))
	}
	var findings []interfaces.Finding
	var srcFiles []string
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
			subDirFindings := RunLinterChecks(subDirPath, checks, depth, currentDepth+1, parallel)
			if len(subDirFindings) > 0 {
				findings = append(findings, subDirFindings...)
			}

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

	// We will pass state containing auxiliary information to the checks, such as function declarations, to avoid redundant parsing and improve performance.
	var wg sync.WaitGroup
	// Wait only for file parsing goroutines
	wg.Add(len(srcFiles))
	state := &interfaces.State{Packages: make(map[string]interfaces.Package), SourceAsts: make(map[string]interfaces.SourceAst)}
	funcDeclsCh := make(chan packages.FuncDeclResult, 10)
	astFileCh := make(chan interfaces.SourceAst, 10)

	// Single consumer for funcDeclsCh — will exit when channel is closed
	var consumerWg sync.WaitGroup
	consumerWg.Add(1)
	go func() {
		defer consumerWg.Done()
		for funcDeclResult := range funcDeclsCh {
			if _, exists := state.Packages[funcDeclResult.PackagePath]; !exists {
				state.Packages[funcDeclResult.PackagePath] = interfaces.Package{FuncDecls: funcDeclResult.FuncDecls}
			}
		}
	}()
	consumerWg.Add(1)
	go func() {
		defer consumerWg.Done()
		for astResult := range astFileCh {
			state.SourceAsts[astResult.Fset.Position(astResult.AstFile.Pos()).Filename] = astResult
		}
	}()

	for _, filePath := range srcFiles {
		// Populate state with source files ASTs funcDecls from imported packages
		go func(filePath string) {
			defer wg.Done()
			astFile, fset, err := ast.GetAst(filePath)
			if err != nil {
				fmt.Printf("Error generating AST for file %s: %s\n", filePath, err)
				return
			}
			// Store the AST and FileSet in the state for later use by checks
			astFileCh <- interfaces.SourceAst{Fset: fset, AstFile: astFile}
			packages.ImportPackages(astFile, funcDeclsCh, state)
		}(filePath)

	}

	// Wait for all producers, then close the channel so the consumer can finish
	wg.Wait()
	close(funcDeclsCh)
	close(astFileCh)
	consumerWg.Wait()

	if parallel {
		findings = append(findings, runChecksInParallel(srcFiles, checks, state)...)
	} else {
		findings = append(findings, runChecksSerially(srcFiles, checks, state)...)
	}

	return findings
}

func runChecksInParallel(srcFiles []string, checks []interfaces.CheckFunc, state *interfaces.State) []interfaces.Finding {

	var findings []interfaces.Finding
	totalJobs := len(srcFiles) * len(checks)
	if totalJobs == 0 {
		return nil
	}

	resultsCh := make(chan []interfaces.Finding, 10)

	for _, filePath := range srcFiles {
		go func(filePath string, state *interfaces.State) {
			astFile, fset := state.SourceAsts[filePath].AstFile, state.SourceAsts[filePath].Fset
			for _, check := range checks {
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

func runChecksSerially(srcFiles []string, checks []interfaces.CheckFunc, state *interfaces.State) []interfaces.Finding {
	var findings []interfaces.Finding
	for _, filePath := range srcFiles {
		astFile, fset := state.SourceAsts[filePath].AstFile, state.SourceAsts[filePath].Fset
		for _, check := range checks {
			res := check(fset, astFile, state)
			if len(res) > 0 {
				findings = append(findings, res...)
			}
		}
	}
	return findings
}

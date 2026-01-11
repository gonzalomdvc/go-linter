package main

import (
	"fmt"
	goast "go/ast"
	"go/token"
	"os"
	"regexp"

	"github.com/gonzalomdvc/go-linter/ast"
	"github.com/gonzalomdvc/go-linter/interfaces"
)

var MaxDepth = 3

func RunLinterChecks(dirname string, checks []func(*token.FileSet, *goast.File) []interfaces.Finding, depth int, currentDepth int) []interfaces.Finding {
	files, err := os.ReadDir(dirname)
	if err != nil {
		panic(fmt.Sprintf("Error reading source code files: %s", err))
	}
	var findings []interfaces.Finding
	var srcFiles []string
	for _, file := range files {
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
			subDirFindings := RunLinterChecks(subDirPath, checks, depth, currentDepth+1)
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

	totalJobs := len(srcFiles) * len(checks)
	if totalJobs == 0 {
		return nil
	}

	resultsCh := make(chan []interfaces.Finding, totalJobs)

	for _, filePath := range srcFiles {
		astFile, fset, err := ast.GetAst(filePath)
		if err != nil {
			panic(fmt.Sprintf("Error generating AST for file %s: %s", filePath, err))
		}
		for _, check := range checks {
			go func(fset *token.FileSet, af *goast.File, chk func(*token.FileSet, *goast.File) []interfaces.Finding) {
				res := chk(fset, af)
				resultsCh <- res
			}(fset, astFile, check)
		}
	}

	for i := 0; i < totalJobs; i++ {
		res := <-resultsCh
		findings = append(findings, res...)
	}
	close(resultsCh)
	return findings
}

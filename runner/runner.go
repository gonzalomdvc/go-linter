package runner

import (
	"fmt"
	goast "go/ast"
	"go/token"
	"os"
	"regexp"

	"github.com/gonzalomdvc/go-linter/ast"
	"github.com/gonzalomdvc/go-linter/checks"
	"github.com/gonzalomdvc/go-linter/interfaces"
)

var MaxDepth = 3

var Checks = []func(*token.FileSet, *goast.File) []interfaces.Finding{
	checks.GL1,
	checks.GL2,
	checks.GL3,
	checks.GL4,
	checks.GL5,
	checks.GL6,
	checks.GL7,
	checks.GL8,
}

func RunLinterChecks(dirname string, checks []func(*token.FileSet, *goast.File) []interfaces.Finding, depth int, currentDepth int, parallel bool) []interfaces.Finding {
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

	if parallel {
		findings = append(findings, runChecksInParallel(srcFiles, checks)...)
	} else {
		findings = append(findings, runChecksSerially(srcFiles, checks)...)
	}

	return findings
}

func runChecksInParallel(srcFiles []string, checks []func(*token.FileSet, *goast.File) []interfaces.Finding) []interfaces.Finding {
	var findings []interfaces.Finding
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

func runChecksSerially(srcFiles []string, checks []func(*token.FileSet, *goast.File) []interfaces.Finding) []interfaces.Finding {
	var findings []interfaces.Finding
	for _, filePath := range srcFiles {
		astFile, fset, err := ast.GetAst(filePath)
		if err != nil {
			panic(fmt.Sprintf("Error generating AST for file %s: %s", filePath, err))
		}
		for _, check := range checks {
			res := check(fset, astFile)
			if len(res) > 0 {
				findings = append(findings, res...)
			}
		}
	}
	return findings
}

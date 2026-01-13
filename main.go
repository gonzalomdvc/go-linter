package main

import (
	"os"

	"github.com/gonzalomdvc/go-linter/runner"
	"github.com/gonzalomdvc/go-linter/ui"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	findings := runner.RunLinterChecks(dir, runner.Checks, 3, 0)
	if len(findings) > 0 {
		ui.PrintFindings(findings)
		os.Exit(1)
	} else {
		ui.PrintSuccessfulMessage()
		os.Exit(0)
	}
}

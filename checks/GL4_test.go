package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/packages"
)

func Test_GL4(t *testing.T) {
	positions := []Position{
		{
			Column: 5,
			Line:   4,
		},
	}
	err := RunCheckTest("GL4.go", true, positions, GL4, &packages.State{})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

}

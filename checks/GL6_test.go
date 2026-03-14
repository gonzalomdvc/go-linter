package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/packages"
)

func Test_GL6(t *testing.T) {
	positions := []Position{
		{
			Column: 2,
			Line:   4,
		},
	}
	err := RunCheckTest("GL6.go", true, positions, GL6, &packages.State{})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

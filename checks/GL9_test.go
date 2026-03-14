package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/packages"
)

func Test_GL9(t *testing.T) {
	positions := []Position{
		{
			Column: 2,
			Line:   6,
		},
	}
	err := RunCheckTest("GL9.go", true, positions, GL9, &packages.State{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/packages"
)

func Test_GL12(t *testing.T) {
	positions := []Position{
		{
			Line:   6,
			Column: 2,
		},
	}
	err := RunCheckTest("GL12.go", true, positions, GL12, &packages.State{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/packages"
)

func Test_GL13(t *testing.T) {
	positions := []Position{
		{
			Line:   7,
			Column: 2,
		},
	}
	err := RunCheckTest("GL13.go", true, positions, GL13, &packages.State{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

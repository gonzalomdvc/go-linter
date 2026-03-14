package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/packages"
)

func Test_GL2(t *testing.T) {
	positions := []Position{
		{
			Column: 2,
			Line:   4,
		},
		{
			Column: 9,
			Line:   6,
		},
	}

	err := RunCheckTest("GL2.go", true, positions, GL2, &packages.State{})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

}

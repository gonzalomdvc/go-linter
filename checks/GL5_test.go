package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/packages"
)

func Test_GL5(t *testing.T) {
	positions := []Position{
		{
			Column: 2,
			Line:   8,
		},
	}
	err := RunCheckTest("GL5.go", true, positions, GL5, &packages.State{})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

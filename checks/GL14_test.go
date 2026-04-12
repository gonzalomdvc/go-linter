package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/packages"
)

func Test_GL14(t *testing.T) {
	positions := []Position{
		{
			Line:   4,
			Column: 2,
		},
		{
			Line:   5,
			Column: 2,
		},
	}
	err := RunCheckTest("GL14.go", true, positions, GL14, &packages.State{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

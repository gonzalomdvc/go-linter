package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/packages"
)

func Test_GL8(t *testing.T) {
	positions := []Position{
		{
			Column: 3,
			Line:   7,
		},
	}
	err := RunCheckTest("GL8.go", true, positions, GL8, &packages.State{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

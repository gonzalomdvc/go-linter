package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/packages"
)

func Test_GL7(t *testing.T) {
	positions := []Position{
		{
			Column: 1,
			Line:   7,
		},
	}
	err := RunCheckTest("GL7.go", true, positions, GL7, &packages.State{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

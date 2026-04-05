package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/packages"
)

func Test_GL11(t *testing.T) {
	positions := []Position{
		{
			Column: 1,
			Line:   3,
		},
	}
	err := RunCheckTest("GL11.go", true, positions, GL11, &packages.State{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

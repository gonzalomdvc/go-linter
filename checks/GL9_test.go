package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/interfaces"
	"github.com/gonzalomdvc/go-linter/test"
)

func Test_GL9(t *testing.T) {
	positions := []interfaces.Position{
		{
			Column: 2,
			Line:   6,
		},
	}
	err := test.RunCheckTest("GL9.go", true, positions, GL9)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

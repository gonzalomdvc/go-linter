package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/interfaces"
	"github.com/gonzalomdvc/go-linter/test"
)

func Test_GL6(t *testing.T) {
	positions := []interfaces.Position{
		{
			Column: 2,
			Line:   4,
		},
	}
	err := test.RunCheckTest("GL6.go", true, positions, GL6)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/interfaces"
	"github.com/gonzalomdvc/go-linter/test"
)

func Test_GL8(t *testing.T) {
	positions := []interfaces.Position{
		{
			Column: 3,
			Line:   7,
		},
	}
	err := test.RunCheckTest("GL8.go", true, positions, GL8)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

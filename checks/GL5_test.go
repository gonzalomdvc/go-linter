package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/interfaces"
	"github.com/gonzalomdvc/go-linter/test"
)

func Test_GL5(t *testing.T) {
	positions := []interfaces.Position{
		{
			Column: 2,
			Line:   8,
		},
	}
	err := test.RunCheckTest("GL5.go", true, positions, GL5)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

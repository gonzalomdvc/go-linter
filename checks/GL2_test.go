package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/interfaces"
	"github.com/gonzalomdvc/go-linter/test"
)

func Test_GL2(t *testing.T) {
	positions := []interfaces.Position{
		{
			Column: 2,
			Line:   4,
		},
		{
			Column: 9,
			Line:   6,
		},
	}

	err := test.RunCheckTest("GL2.go", true, positions, GL2)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

}

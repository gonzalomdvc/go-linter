package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/interfaces"
	"github.com/gonzalomdvc/go-linter/test"
)

func Test_GL4(t *testing.T) {
	positions := []interfaces.Position{
		{
			Column: 5,
			Line:   4,
		},
	}
	err := test.RunCheckTest("GL4.go", true, positions, GL4)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

}

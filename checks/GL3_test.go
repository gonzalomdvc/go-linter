package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/interfaces"
	"github.com/gonzalomdvc/go-linter/test"
)

func Test_GL3(t *testing.T) {
	positions := []interfaces.Position{
		{
			Column: 2,
			Line:   4,
		},
	}
	err := test.RunCheckTest("GL3.go", true, positions, GL3)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

}

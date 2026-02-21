package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/interfaces"
	"github.com/gonzalomdvc/go-linter/test"
)

func Test_GL10(t *testing.T) {
	positions := []interfaces.Position{
		{
			Column: 9,
			Line:   8,
		},
	}
	err := test.RunCheckTest("GL10.go", true, positions, GL10)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/interfaces"
	"github.com/gonzalomdvc/go-linter/test"
)

func Test_GL7(t *testing.T) {
	positions := []interfaces.Position{
		{
			Column: 1,
			Line:   7,
		},
	}
	err := test.RunCheckTest("GL7.go", true, positions, GL7)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

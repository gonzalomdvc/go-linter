package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/interfaces"
	"github.com/gonzalomdvc/go-linter/test"
)

func _(t *testing.T) {
	positions := []interfaces.Position{}
	err := test.RunCheckTest("set file here []", true, positions, GLX)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

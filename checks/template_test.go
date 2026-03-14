package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/packages"
)

func _(t *testing.T) {
	positions := []Position{}
	err := RunCheckTest("set file here []", true, positions, GLX, &packages.State{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

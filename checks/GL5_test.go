package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/test"
)

func Test_GL5(t *testing.T) {
	findings, err := test.RunCheckTest("GL5.go", true, GL5)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(findings) != 1 {
		t.Errorf("Expected 1 finding, got %d", len(findings))
	}
}

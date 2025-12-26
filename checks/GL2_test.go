package checks

import (
	"testing"

	"github.com/gonzalomdvc/go-linter/test"
)

func Test_GL2(t *testing.T) {
	findings, err := test.RunCheckTest("GL2.go", true, GL2)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(findings) != 2 {
		t.Errorf("Expected 2 findings, got %d", len(findings))
	}

}

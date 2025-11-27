package ui

import (
	"go/token"
	"testing"

	"github.com/gonzalomdvc/go-linter/interfaces"
)

func Test_PrintPosition(t *testing.T) {
	pos := token.Position{
		Filename: "../interfaces/interfaces.go",
		Line:     10,
		Column:   1,
		Offset:   48,
	}
	check := interfaces.Check{
		Message: "Test error message",
	}
	_, err := PrintPosition(pos, check.Message)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

}
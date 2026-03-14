package ui

import (
	"go/token"
	"testing"
)

func Test_PrintPosition(t *testing.T) {
	pos := token.Position{
		Filename: "../model/types.go",
		Line:     10,
		Column:   1,
		Offset:   48,
	}

	message := "Test error message"
	_, err := PrintPosition(pos, message)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

}

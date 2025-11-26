package ui

import (
	"go/token"
	"testing"
)

func Test_PrintPosition(t *testing.T) {
	pos := token.Position{
		Filename: "interfaces.go",
		Line:     10,
		Column:   1,
		Offset:   48,
	}
	_, err := PrintPosition(pos, 19, "../interfaces")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

}
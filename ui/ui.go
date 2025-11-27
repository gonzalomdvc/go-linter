package ui

import (
	"fmt"
	"go/token"
	"os"
	"strings"
)

const BUFFER_SIZE = 100

func PrintPosition(pos token.Position, message string) (string, error) {
	file, err := os.Open(pos.Filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, BUFFER_SIZE)
	file.ReadAt(buffer, int64(pos.Offset))

	output := fmt.Sprintf("Warning in file: %s, Line: %d, Column: %d\n", pos.Filename, pos.Line, pos.Column)
	output += fmt.Sprintf("%s\n", message)
	output += fmt.Sprintf("--> %s", strings.Split(string(buffer), "\n")[0])

	return output, nil
}

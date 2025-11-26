package ui

import (
	"go/token"
	"os"
)

func PrintPosition(pos token.Position, bufferSize int, directory string) (string, error) {
	file, err := os.Open(directory + "/" + pos.Filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, bufferSize)
	file.ReadAt(buffer, int64(pos.Offset))

	return string(buffer), nil
}
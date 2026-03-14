package model

import "go/token"

type Finding struct {
	Position token.Position
	Message  string
}

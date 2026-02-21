package test

import (
	"github.com/gonzalomdvc/go-linter/test/GL10_helper"
)

func GL10() (int, int) {
	return GL10_helper.DivAndRemainder(2, 1)
}

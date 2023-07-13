package adt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	l := Cons(1, Cons(2, Cons(3, Null[int]())))
	assert.Equal(t, 3, length(l))
}

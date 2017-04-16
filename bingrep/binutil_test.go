package main

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindIndex(t *testing.T) {
	key := []byte{'d', 'e'}
	target := strings.NewReader("abcdefg")
	assert.Equal(t, 3, FindIndex(target, key))
	target.Seek(0, io.SeekStart)
	assert.Equal(t, -1, FindIndex(target, key, BufferSize(1)))
	target.Seek(0, io.SeekStart)
	assert.Equal(t, -1, FindIndex(target, key, BufferSize(2)))
	target.Seek(0, io.SeekStart)
	assert.Equal(t, 3, FindIndex(target, key, BufferSize(3)))
	target.Seek(0, io.SeekStart)
	assert.Equal(t, 3, FindIndex(target, key, BufferSize(5)))
	target.Seek(0, io.SeekStart)
	assert.Equal(t, 3, FindIndex(target, key, BufferSize(100)))
}

package tools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomString(t *testing.T) {
	n := 4
	res := RandomString(n)
	assert.Equal(t, len(res), n)
}

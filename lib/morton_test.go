package lib_test

import (
	"testing"

	"github.com/mikydna/sports/lib"
	"github.com/stretchr/testify/assert"
)

func TestLib_Morton8(t *testing.T) {
	xyz := [3]uint16{1, 2, 3}
	m := lib.Encode3(xyz)
	assert.Equal(t, lib.Decode3(m), xyz)
}

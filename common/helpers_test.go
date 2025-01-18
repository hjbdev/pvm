package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Version_Compare(t *testing.T) {
	v1 := Version{Major: 1, Minor: 2, Patch: 3, ThreadSafe: false}
	v2 := Version{Major: 1, Minor: 2, Patch: 4}
	v3 := Version{Major: 1, Minor: 2, Patch: 3, ThreadSafe: true}

	assert.Equal(t, v1.LessThan(v2), true)
	assert.Equal(t, v1.Same(v3), false)

	// testing versions with "nulled" (-1) values

	v4 := Version{Major: 1, Minor: 2, Patch: -1}
	v5 := Version{Major: 1, Minor: 2, Patch: 3}

	assert.Equal(t, v4.LessThan(v5), false)

	v6 := Version{Major: 1, Minor: -1}
	v7 := Version{Major: 1, Minor: 2}

	assert.Equal(t, v6.LessThan(v7), false)
}

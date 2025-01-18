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
	assert.Equal(t, v1.LessThanOrEqual(v2), true)
	assert.Equal(t, v1.LessThanOrEqual(v3), true)
	assert.Equal(t, v1.GreaterThan(v2), false)
	assert.Equal(t, v1.GreaterThanOrEqual(v2), false)
	assert.Equal(t, v1.GreaterThanOrEqual(v3), true)
	assert.Equal(t, v1.Equal(v2), false)
	assert.Equal(t, v1.Equal(v3), true)
	assert.Equal(t, v1.Same(v3), false)

	// testing versions with "nulled" (-1) values

	v4 := Version{Major: 1, Minor: 2, Patch: -1}
	v5 := Version{Major: 1, Minor: 2, Patch: 3}

	assert.Equal(t, v4.LessThan(v5), false)
	assert.Equal(t, v4.LessThanOrEqual(v5), true)
	assert.Equal(t, v4.GreaterThan(v5), false)
	assert.Equal(t, v4.GreaterThanOrEqual(v5), true)
	assert.Equal(t, v4.Equal(v5), true)

	v6 := Version{Major: 1, Minor: -1}
	v7 := Version{Major: 1, Minor: 2}

	assert.Equal(t, v6.LessThan(v7), false)
	assert.Equal(t, v6.LessThanOrEqual(v7), true)
	assert.Equal(t, v6.GreaterThan(v7), false)
	assert.Equal(t, v6.GreaterThanOrEqual(v7), true)
	assert.Equal(t, v6.Equal(v7), true)
}

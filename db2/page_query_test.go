package db2

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPageQuery(t *testing.T) {
	var p PageQuery
	var err error

	p, err = NewPageQuery("10", "desc", 15)
	require.NoError(t, err)
	assert.Equal(t, "10", p.Cursor)
	assert.Equal(t, "desc", p.Order)
	assert.Equal(t, uint64(15), p.Limit)

	// Defaults
	p, err = NewPageQuery("", "", 0)
	require.NoError(t, err)
	assert.Equal(t, "asc", p.Order)
	c, err := p.CursorInt64()
	require.NoError(t, err)
	assert.Equal(t, int64(0), c)
	assert.Equal(t, uint64(10), p.Limit)
	p, err = NewPageQuery("", "desc", 0)
	require.NoError(t, err)
	c, err = p.CursorInt64()
	require.NoError(t, err)
	assert.Equal(t, int64(9223372036854775807), c)

	// Max
	p, err = NewPageQuery("", "", 200)
	require.NoError(t, err)

	// Error states
	_, err = NewPageQuery("", "foo", 0)
	assert.Error(t, err)
	_, err = NewPageQuery("", "bar", 201)
	assert.Error(t, err)
}

func TestPageQuery_CursorInt64(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var p PageQuery
	var err error

	p = MustPageQuery("1231-4456", "asc", 0)
	l, r, err := p.CursorInt64Pair("-")
	require.NoError(err)
	assert.Equal(int64(1231), l)
	assert.Equal(int64(4456), r)

	// Defaults
	p = MustPageQuery("", "asc", 0)
	l, r, err = p.CursorInt64Pair("-")
	require.NoError(err)
	assert.Equal(int64(0), l)
	assert.Equal(int64(0), r)
	p = MustPageQuery("", "desc", 0)
	l, r, err = p.CursorInt64Pair("-")
	require.NoError(err)
	assert.Equal(int64(math.MaxInt64), l)
	assert.Equal(int64(math.MaxInt64), r)
	p = MustPageQuery("0", "", 0)
	_, r, err = p.CursorInt64Pair("-")
	require.NoError(err)
	assert.Equal(int64(math.MaxInt64), r)

	// Errors
	p = MustPageQuery("123-foo", "", 0)
	_, _, err = p.CursorInt64Pair("-")
	assert.Error(err)
	p = MustPageQuery("foo-123", "", 0)
	_, _, err = p.CursorInt64Pair("-")
	assert.Error(err)
	p = MustPageQuery("-1:123", "", 0)
	_, _, err = p.CursorInt64Pair("-")
	assert.Error(err)
	p = MustPageQuery("111:-123", "", 0)
	_, _, err = p.CursorInt64Pair("-")
	assert.Error(err)

}

package charts

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPoints_Shift(t *testing.T) {
	ts := [3]int64{1, 2, 3}

	points := Points{
		{time.Unix(1, 0), &ts[0]},
		{time.Unix(2, 0), &ts[1]},
		{time.Unix(3, 0), &ts[2]},
	}
	l := len(points)
	points.Shift()
	assert.Len(t, points, l)
	assert.EqualValues(t, (*int64)(nil), points[l-1].Value)
	assert.Equal(t, time.Unix(4, 0), points[l-1].Timestamp)
}

func TestNewPointsPanics(t *testing.T) {
	cases := []struct {
		name   string
		count  uint
		bucket time.Duration
	}{
		{"invalid count", 1, 10},
		{"negative duration", 2, -10},
		{"zero duration", 2, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for i := uint(0); i < 2; i++ {
				assert.Panics(t, func() {
					NewPoints(tc.count, tc.bucket, time.Now())
				})
			}
		})
	}
}

func TestNewPoints(t *testing.T) {
	base := time.Unix(1337, 0)
	points := NewPoints(3, 1*time.Second, base)
	assert.Len(t, points, 3)
	assert.True(t, points[2].Timestamp.After(base))
	assert.Equal(t, points[2].Timestamp, base.Add(1*time.Second))
}

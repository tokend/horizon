package charts

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPoints_Shift(t *testing.T) {
	points := Points{
		{time.Unix(1, 0), 1},
		{time.Unix(2, 0), 2},
		{time.Unix(3, 0), 3},
	}
	l := len(points)
	points.Shift()
	assert.Len(t, points, l)
	assert.EqualValues(t, 0, points[l-1].Value)
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

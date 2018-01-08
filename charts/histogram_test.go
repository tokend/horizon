package charts

import (
	"testing"

	"time"

	"fmt"

	"github.com/stretchr/testify/assert"
)

//check Render for doesn't returns nil values
func TestRender(t *testing.T) {
	inputVal := []int64{1, 1, 0, 0, 2, 1}

	durations := []time.Duration{
		5 * time.Minute, 15 * time.Minute, 25 * time.Minute,
		35 * time.Minute, 45 * time.Minute, 55 * time.Minute,
	}

	p := make(Points, 6)
	for i, d := range durations {
		p[i].Timestamp = time.Now().Add(-d)
		if i != 2 && i != 3 {
			p[i].Value = &inputVal[i]
		}
	}

	h := Histogram{
		duration: time.Hour,
		points:   p,
		preceded: nil,
	}

	h.points = h.Render()
	for i := range p {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			assert.NotNil(t, h.points[i].Value)
		})
	}
}

//no initial value has been set
func TestHistogram_Render(t *testing.T) {

	durations := []time.Duration{
		5 * time.Minute, 15 * time.Minute, 25 * time.Minute,
		35 * time.Minute, 45 * time.Minute, 55 * time.Minute,
	}

	p := make(Points, 6)

	//var input int64 = 1
	for i, d := range durations {
		p[i].Timestamp = time.Now().Add(-d)
	}

	h := Histogram{
		duration: time.Hour,
		points:   p,
		preceded: nil,
	}

	for i := range p {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			assert.NotPanics(t, func() {
				h.points = h.Render()
			})
		})
	}
}

func TestHistogram_Render2(t *testing.T) {
	inputVal := []int64{0, 1, 1, 1, 1, 1}

	durations := []time.Duration{
		5 * time.Minute, 15 * time.Minute, 25 * time.Minute,
		35 * time.Minute, 45 * time.Minute, 55 * time.Minute,
	}

	points := make(Points, 6)
	for i, d := range durations {
		points[i].Timestamp = time.Now().Add(-d)
		if i >= 1 {
			points[i].Value = &inputVal[i]
		}
	}

	point := Point{
		Timestamp: time.Now().Add(-durations[0]),
		Value:     &inputVal[1],
	}

	h := Histogram{
		duration: time.Hour,
		points:   points,
		preceded: &point,
	}

	h.points = h.Render()
	t.Run("Test zero index", func(t *testing.T) {
		assert.EqualValues(t, *h.preceded.Value, *h.points[0].Value)
	})
}

func TestHistogram_Run(t *testing.T) {
	durations := []time.Duration{
		5 * time.Minute, 15 * time.Minute, 25 * time.Minute,
		35 * time.Minute, 45 * time.Minute, 55 * time.Minute,
	}

	points := make(Points, 6)
	for i, d := range durations {
		points[i].Timestamp = time.Now().Add(-d)
	}

	h := Histogram{
		duration: time.Hour,
		points:   points,
		preceded: nil,
	}
	fmt.Println(h)

	h.Run(10, time.Now().Add(-time.Minute))

	fmt.Println(h)
}

func TestPoints_Shift(t *testing.T) {
	ts := [3]int64{1, 2, 3}

	points := Points{
		{time.Unix(1, 0), &ts[0]},
		{time.Unix(2, 0), &ts[1]},
		{time.Unix(3, 0), &ts[2]},
	}

	h := Histogram{
		duration: 3 * time.Second,
		points:   points,
		preceded: nil,
	}
	l := len(h.points)

	h.Shift()
	assert.Len(t, h.points, l)
	assert.Nil(t, h.points[l-1].Value)
	assert.Equal(t, time.Unix(4, 0), points[l-1].Timestamp)
}

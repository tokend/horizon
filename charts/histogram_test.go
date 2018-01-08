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

	for i := range h.points {
		h.Run(int64(i), time.Now().Add(-durations[i]))
		assert.Equal(t, int64(i), *h.points[i].Value)
	}
}

func TestGetIndex(t *testing.T) {
	/*
		1)задать гистрограмму
		2)убедиться:
			*возвращает правильный индекс
			*сэттит значение в средину
			*сэттит значение в nil
	*/
	points := NewPoints(6, time.Hour/time.Duration(6), time.Now().UTC())
	h := Histogram{
		duration: time.Hour,
		points:   points,
	}

	//h.getIndex()
}

func TestNewHistogram(t *testing.T) {
	h := NewHistogram(time.Hour, 60)

	assert.Equal(t, time.Hour, h.duration)
	assert.Equal(t, 60, len(h.points))
}

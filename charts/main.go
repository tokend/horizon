package charts

import (
	"time"

	"gitlab.com/swarmfund/horizon/db2/history"
)

type point struct {
	time  time.Time
	value uint64
}

type Histogram struct {
	Interval      time.Duration
	Count         uint64
	MaxBorderTime time.Time
	timeInterval  []point
}

func NewHistogram(duration time.Duration, count uint64) *Histogram {
	h := Histogram{
		Interval:      duration,
		Count:         count,
		MaxBorderTime: time.Now().UTC(),
	}

	h.timeInterval = make([]point, h.Count)
	for i := time.Duration(count); i > 0; i-- {
		h.timeInterval[i-1].time = h.MaxBorderTime.Add(-1 * (i - 1) * h.bucketLength())
	}

	go h.Ticker()

	return &h
}

func (h *Histogram) bucketLength() time.Duration {
	return h.Interval / time.Duration(h.Count)
}

func (h *Histogram) shift() {
	h.timeInterval = h.timeInterval[1:]
	h.timeInterval = append(h.timeInterval, point{
		time: time.Now().UTC(),
	})
}

func (h *Histogram) Ticker() {
	ticker := time.NewTicker(h.bucketLength())
	for ; ; <-ticker.C {
		h.shift()
		h.MaxBorderTime.Add(h.bucketLength())
	}
}

func (h *Histogram) Run(entryValue uint64, txTime time.Time) {
	minT := h.MaxBorderTime.Add(-h.Interval)

	if txTime.After(minT) && txTime.Before(h.MaxBorderTime) {

		n := int64(h.Interval) / int64(h.Count)
		h.insert(entryValue, txTime, minT, time.Duration(n))
	}
}

func (h *Histogram) insert(
	entryValue uint64,
	txTime, minT time.Time, interval time.Duration,
) {
	insertIndex := txTime.Sub(minT) / interval

	insertData := point{
		time:  txTime,
		value: entryValue,
	}

	h.findAverage(insertIndex, insertData)
}

func (h *Histogram) findAverage(insertIndex time.Duration, insertData point) {

	if !h.timeInterval[insertIndex].time.IsZero() {
		insertData.value = (h.timeInterval[insertIndex].value + insertData.value) / 2
	}

	p := point{
		insertData.time, insertData.value,
	}

	h.timeInterval[insertIndex] = p
}

type TxHistoryStorage struct {
	TxHistory history.TransactionsQI
}

type Point struct {
	Timestamp time.Time `json:"timestamp"`
	Value     int64     `json:"value"`
}

func (h *Histogram) Render() []Point {
	points := make([]Point, 0, h.Count-1)
	for i := 1; i < len(h.timeInterval); i++ {
		value := int64(h.timeInterval[i].value)
		if h.timeInterval[i].value == 0 {
			value = int64(h.timeInterval[i-1].value)
		}
		points = append(points, Point{
			Timestamp: h.timeInterval[i].time,
			Value:     value,
		})
	}
	return points
}

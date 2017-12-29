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
		MaxBorderTime: time.Now(),
	}

	h.timeInterval = make([]point, h.Count)

	return &h
}

func (h *Histogram) shift() {
	h.timeInterval = h.timeInterval[1:]

	var p point
	h.timeInterval = append(h.timeInterval, p)
}

func (h *Histogram) Ticker() {
	intervalPiece := int64(h.Interval) / int64(h.Count)

	ticker := time.NewTicker(time.Duration(intervalPiece))
	go func() {
		for ; ; <-ticker.C {
			h.shift()
			h.MaxBorderTime.Add(time.Duration(intervalPiece))
		}
	}()
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

func (h *Histogram) render() {
	for i := 1; i < len(h.timeInterval); i++ {
		if h.timeInterval[i].value == 0 {
			h.timeInterval[i] = h.timeInterval[i-1]
		}
	}
}

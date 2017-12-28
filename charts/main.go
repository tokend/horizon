package charts

import (
	"time"

	"fmt"

	"gitlab.com/swarmfund/horizon/db2/history"
)

type Histogram struct {
	Interval      time.Duration
	Count         uint64
	MaxBorderTime time.Time
}

func NewHistogram(duration time.Duration, count uint64) *Histogram {
	histogram := Histogram{
		Interval:      duration,
		Count:         count,
		MaxBorderTime: time.Now(),
	}

	timeInterval = make(timeIntervals, histogram.Count)

	return &histogram
}

func (h *Histogram) Shift() {
	timeInterval = timeInterval[1:]

	var p point
	timeInterval = append(timeInterval, p)
}

func Ticker(histogram *Histogram) {
	intervalPiece := int64(histogram.Interval) / int64(histogram.Count)

	ticker := time.NewTicker(time.Duration(intervalPiece))
	go func() {
		for ; ; <-ticker.C {
			histogram.Shift()
			histogram.MaxBorderTime.Add(time.Duration(intervalPiece))
		}
	}()
}

func (histogram *Histogram) Run(entryValue uint64, txTime time.Time) {
	minT := histogram.MaxBorderTime.Add(-histogram.Interval)

	if txTime.After(minT) && txTime.Before(histogram.MaxBorderTime) {

		n := int64(histogram.Interval) / int64(histogram.Count)
		timeInterval.insert(entryValue, txTime, minT, time.Duration(n))
	}
}

type point struct {
	time  time.Time
	value uint64
}

type timeIntervals []point

var timeInterval timeIntervals

func (ti *timeIntervals) insert(
	entryValue uint64,
	txTime, minT time.Time, interval time.Duration,
) {
	insertIndex := txTime.Sub(minT) / interval
	fmt.Println("DIFF = ", int64(txTime.Sub(minT)), int64(interval))

	insertData := point{
		time:  txTime,
		value: entryValue,
	}

	timeInterval.findMiddle(insertIndex, insertData)
}

func (ti *timeIntervals) findMiddle(insertIndex time.Duration, insertData point) {
	if timeInterval[insertIndex].value != 0 {
		insertData.value = (timeInterval[insertIndex].value + insertData.value) / 2
	}
	p := point{
		insertData.time, insertData.value,
	}

	timeInterval[insertIndex] = p
}

type TxHistoryStorage struct {
	TxHistory history.TransactionsQI
}

func render() {
	for i := 1; i < len(timeInterval); i++ {
		if timeInterval[i].value == 0 {
			timeInterval[i] = timeInterval[i-1]
		}
	}
}

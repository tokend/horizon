package charts

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestHistogram_Shift(t *testing.T) {
	var histogram Histogram
	var insertValue [24]point

	//create values
	for i := 0; i < len(insertValue); i++ {
		insertValue[i].value = rand.Uint64()

		randomTime := rand.Int63n(time.Now().Unix()-94608000) + 94608000
		randomNow := time.Unix(randomTime, 0)

		insertValue[i].time = randomNow
	}

	for i := 0; i < len(insertValue); i++ {
		timeInterval = append(timeInterval, insertValue[i])
	}

	fmt.Println("Before: ", timeInterval)

	histogram.Shift()

	fmt.Println("After: ", timeInterval)
}

func TestRender(t *testing.T) {
	var insertValue [10]point

	//insert random time
	for i := 0; i < len(insertValue); i++ {

		randomTime := rand.Int63n(time.Now().Unix()-94608000) + 94608000
		randomNow := time.Unix(randomTime, 0)

		insertValue[i].time = randomNow
	}

	var temp []uint64
	//input values: 0, 1, 0, 0, 3, 0 , 5, 5, 0, 2
	temp = append(temp, 0, 1, 0, 0, 3, 0, 5, 5, 0, 2)

	for i := 0; i < len(insertValue); i++ {
		insertValue[i].value = temp[i]

	}

	for i := 0; i < len(insertValue); i++ {
		timeInterval = append(timeInterval, insertValue[i])
	}

	fmt.Println("Before render():", timeInterval)

	render()
	// output: 0, 1, 1, 1, 3, 3, 5, 5, 5, 2
	fmt.Println("After render():", timeInterval)
}

func TestInsert(t *testing.T) {
	//var histogram Histogram

	hourDuration := time.Hour

	histogram := NewHistogram(hourDuration, 60)

	points := []time.Duration{1 * time.Minute, 2 * time.Hour, 3 * time.Minute}

	for i, point := range points {
		histogram.Run(uint64(i), time.Now().Add(-point))
		fmt.Println(timeInterval)
	}

	//minT := histogram.MaxBorderTime.Add(-histogram.Interval)
	//
	//fmt.Println("Before insert: ", timeInterval)
	//
	////timeInterval.insert(value, insertValue[0].time, minT, histogram.Interval)
	//fmt.Println("After insert: ", timeInterval)
}

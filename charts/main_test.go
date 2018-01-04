package charts

/*
func TestHistogram_Shift(t *testing.T) {
	h := NewHistogram(time.Hour, 6)

	durations := []time.Duration{
		5 * time.Minute, 15 * time.Minute, 25 * time.Minute,
		35 * time.Minute, 45 * time.Minute, 55 * time.Minute,
	}

	values := []uint64{1, 0, 0, 3, 0, 5}

	for i := range durations {
		h.timeInterval[i].value = values[i]
		h.timeInterval[i].time = time.Now().Add(-durations[i])
	}

	fmt.Println("Before: ", h.timeInterval)

	h.shift()
	fmt.Println("After: ", h.timeInterval)
}

func TestRender(t *testing.T) {

	h := NewHistogram(time.Hour, 6)

	durations := []time.Duration{
		5 * time.Minute, 15 * time.Minute, 25 * time.Minute,
		35 * time.Minute, 45 * time.Minute, 55 * time.Minute,
	}

	values := []uint64{1, 0, 0, 3, 0, 5}

	for i := range durations {
		h.timeInterval[i].value = values[i]
		h.timeInterval[i].time = time.Now().Add(-durations[i])
	}

	fmt.Println("Before render ", h.timeInterval)
	h.render()
	fmt.Println("After render ", h.timeInterval)
}

func TestInsert(t *testing.T) {

	hourDuration := time.Hour

	h := NewHistogram(hourDuration, 60)

	points := []time.Duration{1 * time.Minute, 2 * time.Hour,
		3 * time.Minute, 4 * time.Minute, 5 * time.Minute}

	for i, point := range points {
		h.Run(uint64(i), time.Now().Add(-point))
		fmt.Println(h.timeInterval)
	}
}
*/

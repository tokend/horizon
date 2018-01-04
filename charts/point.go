package charts

import "time"

type Point struct {
	Timestamp time.Time
	// Value holds aggregated bucket value,
	// nil implies nothing has been stored yet
	Value *int64
}

package resource

import (
	"time"
)

type Charts map[string]Chart

type Chart []Point

type Point struct {
	Timestamp time.Time `json:"timestamp"`
	Value     string    `json:"value"`
}

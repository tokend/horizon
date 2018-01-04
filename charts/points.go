package charts

import "time"

type Points []Point

func (p *Points) Shift() {
	*p = append((*p)[1:], Point{
		Timestamp: time.Now().UTC(),
	})
}

func (p Points) First() Point {
	return p[0]
}

func (p Points) Last() Point {
	return p[len(p)-1]
}

func (p Points) Insert(idx int64, value int64) {
	if p[idx].Value != 0 {
		value = (p[idx].Value + value) / 2
	}
	p[idx].Value = value
}

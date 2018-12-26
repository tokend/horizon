package generator

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"math"
	"gitlab.com/distributed_lab/logan/v3"
)

type ID struct {
	Significant int32
	seq         uint32
}

func NewID(significant uint32) *ID {
	if significant > math.MaxInt32 {
		panic(errors.From(errors.New("trying to create new generator with significant > MaxInt32"), logan.F{
			"significant": significant,
		}))
	}
	return &ID{
		Significant: int32(significant),
		seq:         0,
	}
}

func NewIDI32(significant int32) *ID {
	if significant < 0 {
		panic(errors.New("negative significant is not supported"))
	}

	return NewID(uint32(significant))
}

func (g *ID) Next() int64 {
	if g.seq == math.MaxUint32 {
		panic(errors.New("failed to generate next ID - overflow of sequence"))
	}

	g.seq++
	return (int64(g.Significant) << 32) | int64(g.seq)
}

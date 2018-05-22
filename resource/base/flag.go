package base

import "github.com/cheekybits/genny/generic"

// Flag - represents one value of binary mask
type Flag struct {
	Name  string `json:"name,omitempty"`
	Value int32  `json:"value"`
}

type flagValueType generic.Number
func (f flagValueType) String() string {
	return ""
}

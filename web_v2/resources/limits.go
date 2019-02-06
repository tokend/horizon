package resources

import (
	"github.com/spf13/cast"
	"gitlab.com/tokend/regources/v2"
)

func NewLimits(limitsID uint64) *regources.Limits {
	return &regources.Limits{
		Key: NewLimitsKey(limitsID),
	}
}

func NewLimitsKey(limitsID uint64) regources.Key {
	return regources.Key{
		ID:   cast.ToString(limitsID),
		Type: regources.TypeLimits,
	}
}

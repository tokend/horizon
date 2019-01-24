package resources

import (
	"strconv"

	"gitlab.com/tokend/regources/v2"
)

//NewParticipantEffectKey - returns new key for ID
func NewParticipantEffectKey(id int64) regources.Key {
	return regources.Key{
		Type: regources.TypeParticipantEffects,
		ID:   strconv.FormatInt(id, 10),
	}
}

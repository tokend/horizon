package history2

// EffectType describe the effect of some operation to the account or particular balance
type EffectType int64


// Effect stores the details of the effect in union switch form. Only one value should be selected
// Effect should never store more than one change to the account or balance
type Effect struct {
	Type EffectType `json:"type"`

}

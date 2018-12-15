package history2

// EffectType describe the effect of some operation to the account or particular balance
type EffectType int64

const (
	EffectTypeNone EffectType = iota
	EffectTypeFunded
	EffectTypeIssued
	EffectTypeCharged
	EffectTypeWithdrawn
	EffectTypeLocked
	EffectTypeUnlocked
	EffectTypeChargedFromLocked
	EffectTypeFundedToLocked
	EffectTypeMatched // funded, charged or charge from locked in offer cases
)

// Effect stores the details of the effect in union switch form. Only one value should be selected
// Effect should never store more than one change to the account or balance
type Effect struct {
	Type              EffectType               `json:"type"`
	Funded            *FundedEffect            `json:"funded,omitempty"`
	Issued            *FundedEffect            `json:"issued,omitempty"`
	Charged           *ChargedEffect           `json:"charged,omitempty"`
	Withdrawn         *ChargedFromLockedEffect `json:"withdrawn,omitempty"`
	Locked            *LockedEffect            `json:"locked,omitempty"`
	Unlocked          *UnlockedEffect          `json:"unlocked,omitempty"`
	ChargedFromLocked *ChargedFromLockedEffect `json:"charged_from_locked,omitempty"`
	FundedToLocked    *FundedToLockedEffect    `json:"funded_to_locked,omitempty"`
	Offer             *OfferEffect             `json:"offer,omitempty"`
	DeletedOffer      *DeletedOfferEffect      `json:"deleted_offer,omitempty"`
}

type OfferEffect struct {
	// maybe add offer id
	BaseBalanceAddress  string  `json:"base_balance_address"`
	QuoteBalanceAddress string  `json:"quote_balance_address"`
	BaseAmount          string  `json:"base_amount"`
	QuoteAmount         string  `json:"quote_amount"`
	BaseAsset           string  `json:"base_asset"`
	QuoteAsset          string  `json:"quote_asset"`
	Price               string  `json:"price"`
	IsBuy               bool    `json:"is_buy"`
	FeePaid             FeePaid `json:"fee_paid"`
}

type DeletedOfferEffect struct {
	BaseAmount string `json:"base_amount"`
}

type FundedEffect struct {
	Amount  string  `json:"amount"`
	FeePaid FeePaid `json:"fee_paid"`
}

type FeePaid struct {
	Fixed             string `json:"fixed"`
	CalculatedPercent string `json:"calculated_percent"`
}

type LockedEffect struct {
	Amount    string  `json:"amount"`
	FeeLocked FeePaid `json:"fee_locked"`
}

type ChargedEffect struct {
	Amount  string  `json:"amount"`
	FeePaid FeePaid `json:"fee_paid"`
}

type UnlockedEffect struct {
	Amount      string  `json:"amount"`
	FeeUnlocked FeePaid `json:"fee_unlocked"`
}

type ChargedFromLockedEffect struct {
	Amount  string  `json:"amount"`
	FeePaid FeePaid `json:"fee_paid"`
}

type FundedToLockedEffect struct {
	Amount  string  `json:"amount"`
	FeePaid FeePaid `json:"fee_paid"`
}

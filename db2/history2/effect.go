package history2

import "gitlab.com/tokend/go/xdr"

// EffectType describe the effect of some operation to the account or particular balance
type EffectType int64

const (
	EffectTypeNone EffectType = iota
	EffectTypeFunded
	EffectTypeCharged
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
	Charged           *ChargedEffect           `json:"withdraw,omitempty"`
	Locked            *LockedEffect            `json:"locked,omitempty"`
	Unlocked          *UnlockedEffect          `json:"unlocked,omitempty"`
	ChargedFromLocked *ChargedFromLockedEffect `json:"charged_from_locked,omitempty"`
	FundedToLocked    *FundedToLockedEffect    `json:"funded_to_locked,omitempty"`
	Offer             *OfferEffect             `json:"offer,omitempty"`
	DeletedOffer      *DeletedOfferEffect      `json:"deleted_offer,omitempty"`
}

type OfferEffect struct {
	// maybe add offer id
	BaseBalanceID  int64         `json:"base_balance_id"`
	QuoteBalanceID int64         `json:"quote_balance_id"`
	BaseAmount     string        `json:"base_amount"`
	QuoteAmount    string        `json:"quote_amount"`
	BaseAsset      xdr.AssetCode `json:"base_asset"`
	QuoteAsset     xdr.AssetCode `json:"quote_asset"`
	Price          string        `json:"price"`
	IsBuy          bool          `json:"is_buy"`
	FeePaid        FeePaid       `json:"fee_paid"`
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

package history2

import "gitlab.com/tokend/go/xdr"

// EffectType describe the effect of some operation to the account or particular balance
type EffectType int64

const (
	EffectTypeNone EffectType = iota
	EffectTypeWithdraw
	EffectTypeOffer
)

// Effect stores the details of the effect in union switch form. Only one value should be selected
// Effect should never store more than one change to the account or balance
type Effect struct {
	Type           EffectType   `json:"type"`
	WithdrawAmount *int64       `json:"withdraw_amount,omitempty"`
	Offer          *OfferEffect `json:"offer_effect,omitempty"`
}

type OfferEffect struct {
	// maybe add offer id
	BaseBalanceID  int64         `json:"base_balance_id"`
	QuoteBalanceID int64         `json:"quote_balance_id"`
	BaseAmount     string        `json:"base_amount"`
	BaseAsset      xdr.AssetCode `json:"base_asset"`
	QuoteAsset     xdr.AssetCode `json:"quote_asset"`
}

// for match for one effect speicify  incommig and outgoing balance and asset and amount in one struct

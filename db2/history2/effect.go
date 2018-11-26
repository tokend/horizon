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
	EffectTypeMatched // funded, charged or charge from locked in offer cases
)

// Effect stores the details of the effect in union switch form. Only one value should be selected
// Effect should never store more than one change to the account or balance
type Effect struct {
	Type     EffectType      `json:"type"`
	Issuance *IssuanceEffect `json:"issuance,omitempty"`
	Withdraw *WithdrawEffect `json:"withdraw,omitempty"`
	Offer    *OfferEffect    `json:"offer,omitempty"`
	AMLAlert *AMLAlertEffect `json:"aml_alert,omitempty"`
	Payment  *PaymentEffect  `json:"payment"`
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
}

type IssuanceEffect struct {
	Amount int64 `json:"amount"`
}

type WithdrawEffect struct {
	Amount int64 `json:"amount"`
}

type AMLAlertEffect struct {
	Amount string `json:"amount"`
}

type PaymentEffect struct {
	Amount string `json:"amount"`
}

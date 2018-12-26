package history2

// effectType describe the effect of some operation to the account or particular balance
type effectType int64

const (
	// EffectTypeFunded - balance received funds from other balance
	EffectTypeFunded effectType = iota
	// EffectTypeIssued - funds have been issued to the balance
	EffectTypeIssued
	// EffectTypeCharged - balance has been charged
	EffectTypeCharged
	// EffectTypeWithdrawn - balance has been charged and corresponding amount of tokens has been destroyed
	EffectTypeWithdrawn
	// EffectTypeLocked - funds has been locked on the balance
	EffectTypeLocked
	// EffectTypeUnlocked - funds has been unlocked on the balance
	EffectTypeUnlocked
	// EffectTypeChargedFromLocked - funds has been charged from locked amount on balance
	EffectTypeChargedFromLocked
	// EffectTypeMatched - balance has been charged or received funds due to match of the offers
	EffectTypeMatched
)

// Effect stores the details of the operation effect on balance in union switch form. Only one value should be selected
// Effect should never store more than one change to the account or balance
type Effect struct {
	Type              effectType           `json:"type"`
	Funded            *BalanceChangeEffect `json:"funded,omitempty"`
	Issued            *BalanceChangeEffect `json:"issued,omitempty"`
	Charged           *BalanceChangeEffect `json:"charged,omitempty"`
	Withdrawn         *BalanceChangeEffect `json:"withdrawn,omitempty"`
	Locked            *BalanceChangeEffect `json:"locked,omitempty"`
	Unlocked          *BalanceChangeEffect `json:"unlocked,omitempty"`
	ChargedFromLocked *BalanceChangeEffect `json:"charged_from_locked,omitempty"`
	Matched           *MatchEffect         `json:"matched"`
}

// MatchEffect - describes changes to base and quote balance occurred on match
type MatchEffect struct {
	OfferID     int64                         `json:"offer_id"`
	OrderBookID int64                         `json:"order_book_id"`
	Price       string                        `json:"price"`
	Charged     ParticularBalanceChangeEffect `json:"charged"`
	Funded      ParticularBalanceChangeEffect `json:"funded"`
}

// ParticularBalanceChangeEffect - describes movement of fund for particular balance
type ParticularBalanceChangeEffect struct {
	BalanceChangeEffect
	BalanceAddress string `json:"balance_address"`
	AssetCode      string `json:"asset_code"`
}

// BalanceChangeEffect - describes movement of funds
type BalanceChangeEffect struct {
	Amount string `json:"amount"`
	Fee    Fee    `json:"fee"`
}

// Fee - describes fee happened on balance. Direction of fee depends on the operation (depending on effect might be charged, locked, unlocked,
// for all incoming effects but unlocked it's always charged)
type Fee struct {
	Fixed             string `json:"fixed"`
	CalculatedPercent string `json:"calculated_percent"`
}

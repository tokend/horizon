package rgenerated

type QuoteAssetAttributes struct {
	CurrentCap      Amount  `json:"current_cap"`
	HardCap         Amount  `json:"hard_cap"`
	Price           Amount  `json:"price"`
	SoftCap         *Amount `json:"soft_cap,omitempty"`
	TotalCurrentCap Amount  `json:"total_current_cap"`
}

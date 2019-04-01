package rgenerated

type EffectMatchedAttributes struct {
	Charged     ParticularBalanceChangeEffect `json:"charged"`
	Funded      ParticularBalanceChangeEffect `json:"funded"`
	OfferId     int64                         `json:"offer_id"`
	OrderBookId int64                         `json:"order_book_id"`
	Price       Amount                        `json:"price"`
}

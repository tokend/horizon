package regources

//ManageOffer - details of corresponding op
type ManageOffer struct {
	Key
	Attributes ManageOfferAttrs `json:"attributes"`
}

//ManageOfferAttrs - details of corresponding op
type ManageOfferAttrs struct {
	OfferID     int64  `json:"offer_id,omitempty"`
	OrderBookID int64  `json:"order_book_id"`
	BaseAsset   string `json:"base_asset"`
	QuoteAsset  string `json:"quote_asset"`
	Amount      Amount `json:"base_amount"`
	Price       Amount `json:"price"`
	IsBuy       bool   `json:"is_buy"`
	Fee         Fee    `json:"fee"`
	IsDeleted   bool   `json:"is_deleted"`
}


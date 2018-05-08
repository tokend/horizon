package operations

type PaymentV2 struct {
	Base
	PaymentID           uint64                 `json:"payment_id"`
	From                string                 `json:"from,omitempty"`
	To                  string                 `json:"to,omitempty"`
	FromBalance         string                 `json:"from_balance,omitempty"`
	ToBalance           string                 `json:"to_balance,omitempty"`
	Amount              string                 `json:"amount"`
	Asset               string                 `json:"asset"`
	SourceFeeData       map[string]interface{} `json:"source_fee_data"`
	DestinationFeeData  map[string]interface{} `json:"destination_fee_data"`
	SourcePaysForDest   bool                   `json:"source_pays_for_dest"`
	Subject             string                 `json:"subject"`
	Reference           string                 `json:"reference"`
	SourceSentUniversal string                 `json:"source_sent_universal"`
}

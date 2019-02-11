package regources

// WithdrawalRequest - represents details of the `withdrawal` reviewable request
type WithdrawalRequest struct {
	Key
	Attributes    WithdrawalRequestAttrs     `json:"attributes"`
	Relationships WithdrawalRequestRelations `json:"relationships"`
}

// WithdrawalRequestAttrs - attributes of the `withdrawal` reviewable request
type WithdrawalRequestAttrs struct {
	Fee             Fee     `json:"fee"`
	Amount          string  `json:"amount"`
	Details         Details `json:"external_details"`
	ReviewerDetails Details `json:"reviewer_details"`
}

// WithdrawalRequestRelations - relationships of the `withdrawal` reviewable request
type WithdrawalRequestRelations struct {
	Balance *Relation `json:"balance"`
}

package regources

// PaymentAttrs - stores details of payment operation
type Payment struct {
	Key
	Attributes PaymentAttrs `json:"attributes"`
}

// PaymentAttrs - stores details of payment operation
type PaymentAttrs struct {
	AccountFrom             string `json:"account_from"`
	AccountTo               string `json:"account_to"`
	BalanceFrom             string `json:"balance_from"`
	BalanceTo               string `json:"balance_to"`
	Amount                  Amount `json:"amount"`
	Asset                   string `json:"asset"`
	SourceFee               Fee    `json:"source_fee"`
	DestinationFee          Fee    `json:"destination_fee"`
	SourcePayForDestination bool   `json:"source_pay_for_destination"`
	Subject                 string `json:"subject"`
	Reference               string `json:"reference"`
}

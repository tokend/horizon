package regources


//CreateWithdrawRequest - details of corresponding op
type CreateWithdrawRequest struct {
	Key
	Attributes CreateWithdrawRequestAttrs `json:"attributes"`
}

//CreateWithdrawRequestAttrs - details of corresponding op
type CreateWithdrawRequestAttrs struct {
	BalanceAddress  string  `json:"balance_address"`
	Amount          Amount  `json:"amount"`
	Fee             Fee     `json:"fee"`
	ExternalDetails Details `json:"external_details"`
}

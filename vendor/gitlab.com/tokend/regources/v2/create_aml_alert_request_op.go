package regources

//CreateAMLAlertRequest - details of corresponding op
type CreateAMLAlertRequest struct {
	Key
	Attributes CreateAMLAlertRequestAttrs `json:"attributes"`
}

//CreateAMLAlertRequestAttrs - details of corresponding op
type CreateAMLAlertRequestAttrs struct {
	Amount         Amount `json:"amount"`
	BalanceAddress string `json:"balance_address"`
	Reason         string `json:"reason"`
}

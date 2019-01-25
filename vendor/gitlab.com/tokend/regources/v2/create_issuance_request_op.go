package regources

//CreateIssuanceRequestAttrs - details of corresponding op
type CreateIssuanceRequest struct {
	Key
	Attributes CreateIssuanceRequestAttrs `json:"attributes"`
}

//CreateIssuanceRequestAttrs - details of corresponding op
type CreateIssuanceRequestAttrs struct {
	Fee                    Fee     `json:"fee"`
	Reference              string  `json:"reference"`
	Amount                 Amount  `json:"amount"`
	Asset                  string  `json:"asset"`
	ReceiverAccountAddress string  `json:"receiver_account_address"`
	ReceiverBalanceAddress string  `json:"receiver_balance_address"`
	ExternalDetails        Details `json:"external_details"`
	AllTasks               *int64  `json:"all_tasks,omitempty"`
	RequestDetails         Request `json:"request_details"`
}

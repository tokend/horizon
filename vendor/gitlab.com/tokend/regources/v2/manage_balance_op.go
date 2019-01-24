package regources

import "gitlab.com/tokend/go/xdr"

//ManageBalance - stores details of create account operation
type ManageBalance struct {
	Key
	Attributes ManageBalanceAttrs `json:"attributes"`
}

//ManageBalanceAttrs - details of ManageBalanceOp
type ManageBalanceAttrs struct {
	DestinationAccount string                  `json:"destination_account"`
	Action             xdr.ManageBalanceAction `json:"action"`
	Asset              string                  `json:"asset"`
	BalanceAddress     string                  `json:"balance_address"`
}

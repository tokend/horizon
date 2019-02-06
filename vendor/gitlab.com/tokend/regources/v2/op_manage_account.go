package regources

import "gitlab.com/tokend/go/xdr"

//ManageAccount - details of ManageAccountOp
type ManageAccount struct {
	Key
	Attributes ManageAccountAttrs `json:"attributes"`
}

//ManageAccountAttrs - details of ManageAccountOp
type ManageAccountAttrs struct {
	AccountAddress       string           `json:"account_address"`
	BlockReasonsToAdd    xdr.BlockReasons `json:"block_reasons_to_add"`
	BlockReasonsToRemove xdr.BlockReasons `json:"block_reasons_to_remove"`
}

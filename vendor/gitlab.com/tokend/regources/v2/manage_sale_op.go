package regources

import (
	"gitlab.com/tokend/go/xdr"
)

//ManageSale - details of corresponding op
type ManageSale struct {
	Key
	Attributes ManageSaleAttrs `json:"attributes"`
}

//ManageSaleAttrs - details of corresponding op
type ManageSaleAttrs struct {
	SaleID uint64               `json:"sale_id"`
	Action xdr.ManageSaleAction `json:"action"`
}

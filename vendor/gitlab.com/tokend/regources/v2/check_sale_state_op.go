package regources

import "gitlab.com/tokend/go/xdr"

//CheckSaleStateAttrs - details of corresponding op
type CheckSaleState struct {
	Key
	Attributes CheckSaleStateAttrs `json:"attributes"`
}

//CheckSaleStateAttrs - details of corresponding op
type CheckSaleStateAttrs struct {
	SaleID int64                    `json:"sale_id"`
	Effect xdr.CheckSaleStateEffect `json:"effect"`
}

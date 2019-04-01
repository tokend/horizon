package rgenerated

type ManageSaleOpAttributes struct {
	// * 1: \"create_update_details_request\" * 2: \"cancel\"
	Action xdr.ManageSaleAction `json:"action"`
	// ID of the sale to manage
	SaleId uint64 `json:"sale_id"`
}

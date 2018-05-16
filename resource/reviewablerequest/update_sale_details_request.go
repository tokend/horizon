package reviewablerequest

import "gitlab.com/swarmfund/horizon/db2/history"

type UpdateSaleDetailsRequest struct {
	SaleID     uint64                 `json:"sale_id"`
	NewDetails map[string]interface{} `json:"new_details"`
}

func (r *UpdateSaleDetailsRequest) Populate(histRequest history.UpdateSaleDetailsRequest) error {
	r.SaleID = histRequest.SaleID
	r.NewDetails = histRequest.NewDetails
	return nil
}

package reviewablerequest

import "gitlab.com/swarmfund/horizon/db2/history"

type UpdateSaleEndTimeRequest struct {
	SaleID     uint64 `json:"sale_id"`
	NewEndTime uint64 `json:"new_end_time"`
}

func (r *UpdateSaleEndTimeRequest) Populate(histRequest history.UpdateSaleEndTimeRequest) error {
	r.SaleID = histRequest.SaleID
	r.NewEndTime = histRequest.NewEndTime
	return nil
}

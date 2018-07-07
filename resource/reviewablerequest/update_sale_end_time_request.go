package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"time"
)

type UpdateSaleEndTimeRequest struct {
	SaleID     uint64    `json:"sale_id"`
	NewEndTime time.Time `json:"new_end_time"`
}

func (r *UpdateSaleEndTimeRequest) Populate(histRequest history.UpdateSaleEndTimeRequest) error {
	r.SaleID = histRequest.SaleID
	r.NewEndTime = histRequest.NewEndTime
	return nil
}

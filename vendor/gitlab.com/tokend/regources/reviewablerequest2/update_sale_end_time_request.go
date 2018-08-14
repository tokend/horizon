package reviewablerequest2

import "time"

type UpdateSaleEndTimeRequest struct {
	SaleID     uint64    `json:"sale_id"`
	NewEndTime time.Time `json:"new_end_time"`
}

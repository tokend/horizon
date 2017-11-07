package resource

import (
	"gitlab.com/distributed_lab/tokend/horizon/db2/history"
	"gitlab.com/distributed_lab/tokend/horizon/resource/operations"
	"strconv"
	"time"
)

type PaymentRequest struct {
	PT             string                 `json:"paging_token"`
	Exchange       string                 `json:"exchange"`
	PaymentID      string                 `json:"payment_id"`
	PaymentState   uint32                 `json:"payment_state"`
	Accepted       *bool                  `json:"accepted"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	PaymentDetails operations.BasePayment `json:"payment_details"`
	FromEmail      string                 `json:"from_email,omitempty"`
	ToEmail        string                 `json:"to_email,omitempty"`
	RequestType    int32                  `json:"request_type"`
}

// Populate fills out the resource's fields
func (request *PaymentRequest) Populate(row *history.PaymentRequest) error {
	request.PT = strconv.FormatInt(row.ID, 10)
	request.Exchange = row.Exchange
	request.PaymentID = strconv.FormatUint(row.PaymentID, 10)
	if row.PaymentState != nil {
		request.PaymentState = *row.PaymentState
	} else {
		if row.Accepted == nil {
			request.PaymentState = history.PENDING
		} else if *row.Accepted {
			request.PaymentState = history.SUCCESS
		} else {
			request.PaymentState = history.REJECTED
		}
	}
	request.Accepted = row.Accepted
	request.CreatedAt = row.CreatedAt
	request.UpdatedAt = row.UpdatedAt
	request.RequestType = row.RequestType
	err := row.UnmarshalDetails(&request.PaymentDetails)
	return err
}

// PagingToken implementation for hal.Pageable
func (request *PaymentRequest) PagingToken() string {
	return request.PT
}

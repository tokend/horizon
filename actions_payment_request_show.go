package horizon

import (
	"database/sql"

	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/resource"
)

// PaymentRequestShowAction returns a payment request based upon the provided
// id
type PaymentRequestShowAction struct {
	Action
	RequestID   uint64
	byPaymentID bool
	Record      history.PaymentRequest
	Resource    resource.PaymentRequest
}

// JSON is a method for actions.JSON
func (action *PaymentRequestShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecords,
		action.loadPage,
	)
	action.Do(func() {
		hal.Render(action.W, action.Resource)
	})
}

func (action *PaymentRequestShowAction) loadParams() {
	action.RequestID = action.GetUInt64("id")
	action.byPaymentID = action.GetBool("by_payment_id")
}

func (action *PaymentRequestShowAction) loadRecords() {
	if action.byPaymentID {
		action.Err = action.HistoryQ().PaymentRequestByPaymentID(&action.Record, action.RequestID)
	} else {
		action.Err = action.HistoryQ().PaymentRequestByID(&action.Record, action.RequestID)
	}

	if action.Err != nil && action.Err != sql.ErrNoRows {
		action.Log.WithError(action.Err).Error("Failed to get Recovery request from db")
	}
}

func (action *PaymentRequestShowAction) loadPage() {
	action.Resource.Populate(&action.Record)
}

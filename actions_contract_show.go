package horizon

import (
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
	"gitlab.com/swarmfund/horizon/resource/reviewablerequest"
	"gitlab.com/tokend/regources"
)

// ContractShowAction represents singe contract by id
// with invoices requests
type ContractShowAction struct {
	Action
	ContractID     int64
	ContractRecord regources.Contract
}

// JSON is a method for actions.JSON
func (action *ContractShowAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.loadRecords,
		action.checkAllowed,
		func() {
			hal.Render(action.W, action.ContractRecord)
		},
	)
}

func (action *ContractShowAction) checkAllowed() {
	action.IsAllowed(action.ContractRecord.Contractor, action.ContractRecord.Customer, action.ContractRecord.Escrow)
}

func (action *ContractShowAction) loadParams() {
	action.ContractID = action.GetInt64("id")
}

func (action *ContractShowAction) loadRecords() {

	contract, err := action.HistoryQ().Contracts().ByID(action.ContractID)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get contract record")
		action.Err = &problem.ServerError
		return
	}

	if contract == nil {
		action.Err = &problem.NotFound
		return
	}

	action.ContractRecord = resource.PopulateContract(*contract)

	invoices, err := action.HistoryQ().ReviewableRequests().ByIDs(contract.Invoices).Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to get invoices records")
		action.Err = &problem.ServerError
		return
	}

	for _, invoice := range invoices {
		res, err := reviewablerequest.PopulateReviewableRequest(&invoice)
		if err != nil {
			action.Log.WithError(err).Error("Failed to populate invoice request")
			action.Err = &problem.ServerError
			return
		}

		action.ContractRecord.Invoices = append(action.ContractRecord.Invoices, *res)
	}
}

package horizon

import (
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
	"gitlab.com/tokend/regources"
)

// ContractIndexAction renders a page of contracts
// filters by startTime, endTime, disputing state,
// contractorID, customerID
type ContractIndexAction struct {
	Action
	PagingParams     db2.PageQuery
	StartTime        *int64
	EndTime          *int64
	Disputing        *bool
	ContractorID     string
	CustomerID       string
	EscrowID         string
	ContractsRecords []regources.Contract
	Page             hal.Page
}

// JSON is a method for actions.JSON
func (action *ContractIndexAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.checkAllowed,
		action.loadParams,
		action.ValidateCursorWithinHistory,
		action.loadRecords,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *ContractIndexAction) checkAllowed() {
	var allowedAccounts []string
	allowedAccounts = append(allowedAccounts, action.ContractorID)
	allowedAccounts = append(allowedAccounts, action.CustomerID)
	allowedAccounts = append(allowedAccounts, action.EscrowID)
	action.IsAllowed(allowedAccounts...)
}

func (action *ContractIndexAction) loadParams() {
	action.ValidateCursorAsDefault()
	action.StartTime = action.GetOptionalInt64("start_time")
	action.EndTime = action.GetOptionalInt64("end_time")
	action.Disputing = action.GetOptionalBool("state")
	action.ContractorID = action.GetString("contractor_id")
	action.CustomerID = action.GetString("customer_id")
	action.EscrowID = action.GetString("escrow_id")
	action.PagingParams = action.GetPageQuery()
}

func (action *ContractIndexAction) loadRecords() {
	q := action.HistoryQ().Contracts()
	if action.StartTime != nil {
		q = q.ByStartTime(*action.StartTime)
	}
	if action.EndTime != nil {
		q = q.ByEndTime(*action.EndTime)
	}
	if action.Disputing != nil {
		q = q.ByDisputeState(*action.Disputing)
	}
	if action.ContractorID != "" {
		q = q.ByContractorID(action.ContractorID)
	}
	if action.CustomerID != "" {
		q = q.ByCustomerID(action.CustomerID)
	}
	if action.EscrowID != "" {
		q = q.ByEscrowID(action.EscrowID)
	}

	historyContracts, err := q.Page(action.PagingParams).Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to get contracts records")
		action.Err = &problem.ServerError
		return
	}

	for _, contract := range historyContracts {
		action.Page.Add(resource.PopulateContract(contract))
	}
}

func (action *ContractIndexAction) loadPage() {
	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}

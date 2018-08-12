package horizon

import (
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/tokend/regources"
	"gitlab.com/swarmfund/horizon/resource"
)

// TransactionV2IndexAction: pages of transactions

// TransactionV2IndexAction renders a page of ledger resources, identified by
// a normal page query, entry type and effects
type ContractIndexAction struct {
	Action
	PagingParams     db2.PageQuery
	StartTime        *int64
	EndTime          *int64
	Disputing            *bool
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
	action.IsAllowed("")
}

func (action *ContractIndexAction) loadParams() {
	action.ValidateCursorAsDefault()
	action.StartTime = action.GetOptionalInt64("start_time")
	action.EndTime = action.GetOptionalInt64("end_time")
	action.Disputing = action.GetOptionalBool("state")
	action.PagingParams = action.getTxPageQuery()
}
// array in object for invoices json
const (
	maxContractPagSize uint64 = 1000
)

func (action *ContractIndexAction) getTxPageQuery() db2.PageQuery {
	pagingParams := action.GetPageQuery()
	limit := action.GetUInt64("limit")
	if limit > maxContractPagSize {
		pagingParams.Limit = maxContractPagSize
	}

	return pagingParams
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

	historyContracts, err := q.Page(action.PagingParams).Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to get contracts records")
		action.Err = &problem.ServerError
		return
	}

	for _, contract := range historyContracts {
		action.ContractsRecords = append(action.ContractsRecords, resource.PopulateContract(contract))
	}
}

func (action *ContractIndexAction) loadPage() {
	for _, record := range action.ContractsRecords {
		action.Page.Add(record)
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}

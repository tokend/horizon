package horizon

import (
	"strconv"

	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/sse"
	"gitlab.com/swarmfund/horizon/resource"
)

type PaymentRequestIndexAction struct {
	Action
	BalanceFilter  string
	AccountFilter  string
	StateFilter    int32
	OnlyForfeits   bool
	PagingParams   db2.PageQuery
	Records        []history.PaymentRequest
	Page           hal.Page
}

// JSON is a method for actions.JSON
func (action *PaymentRequestIndexAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadRecords,
		action.loadPage,
	)
	action.Do(func() {
		hal.Render(action.W, action.Page)
	})
}

func (action *PaymentRequestIndexAction) SSE(stream sse.Stream) {
	action.Setup(
		action.loadParams,
		action.checkAllowed,
	)
	action.Do(
		func() {
			// we will reuse this variable in sse, so re-initializing is required
			action.Records = []history.PaymentRequest{}
		},
		action.loadRecords,
		func() {
			records := action.Records[:]

			for _, record := range records {
				var request resource.PaymentRequest
				request.Populate(&record)

				stream.Send(sse.Event{
					ID:   request.PagingToken(),
					Data: request,
				})

				action.PagingParams.Cursor = request.PagingToken()
			}
		})
}

func (action *PaymentRequestIndexAction) loadParams() {
	action.ValidateCursorAsDefault()
	action.BalanceFilter = action.GetString("target_balance")
	action.AccountFilter = action.GetString("target_account")
	action.StateFilter = action.GetInt32("state")
	action.PagingParams = action.GetPageQuery()
	action.Page.Filters = map[string]string{
		"state":          strconv.FormatInt(int64(action.StateFilter), 10),
		"target_balance": action.BalanceFilter,
		"target_account": action.AccountFilter,
	}

}

func (action *PaymentRequestIndexAction) loadRecords() {
	q := action.HistoryQ()
	requests := q.PaymentRequests()

	if action.AccountFilter != "" {
		requests.ForAccount(action.AccountFilter)
	}

	if action.BalanceFilter != "" {
		requests.ForBalance(action.BalanceFilter)
	}

	if action.OnlyForfeits {
		requests.ForfeitsOnly()
	}

	action.Err = requests.Page(action.PagingParams).Select(&action.Records)
	if action.Err != nil {
		action.Log.WithError(action.Err).Error("Faieled to get payment requests")
		return
	}
}

func (action *PaymentRequestIndexAction) loadPage() {
	for _, record := range action.Records {
		var request resource.PaymentRequest
		request.Populate(&record)
		action.Page.Add(&request)
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}

func (action *PaymentRequestIndexAction) checkAllowed() {
	action.IsAllowed(action.AccountFilter)
}

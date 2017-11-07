package horizon

import (
	"strconv"

	"gitlab.com/distributed_lab/tokend/horizon/db2"
	"gitlab.com/distributed_lab/tokend/horizon/db2/history"
	"gitlab.com/distributed_lab/tokend/horizon/render/hal"
	"gitlab.com/distributed_lab/tokend/horizon/render/problem"
	"gitlab.com/distributed_lab/tokend/horizon/render/sse"
	"gitlab.com/distributed_lab/tokend/horizon/resource"
)

const (
	REQUEST_PENDING = 1 + iota
	REQUEST_APPROVED
	REQUEST_REJECTED
)

// CoinsEmissionRequestIndexAction returns a paged slice of coins emission requests based upon the provided
// filters
type CoinsEmissionRequestIndexAction struct {
	Action
	ExchangeFilter  string
	ReferenceFilter string
	AssetFilter     string
	StateFilter     int32
	PagingParams    db2.PageQuery
	Records         []history.CoinsEmissionRequest
	Page            hal.Page
}

// JSON is a method for actions.JSON
func (action *CoinsEmissionRequestIndexAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadRecords,
		action.loadExchangeNames,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

// SSE is a method for actions.SSE
func (action *CoinsEmissionRequestIndexAction) SSE(stream sse.Stream) {
	action.Setup(
		action.loadParams,
	)
	action.Do(
		func() {
			// we will reuse this variable in sse, so re-initializing is required
			action.Records = []history.CoinsEmissionRequest{}
		},
		action.loadRecords,
		func() {
			records := action.Records[:]
			for i := range records {
				var request resource.CoinsEmissionRequest
				request.Populate(&records[i])

				stream.Send(sse.Event{
					ID:   request.PagingToken(),
					Data: request,
				})

				action.PagingParams.Cursor = request.PagingToken()
			}
		})
}

func (action *CoinsEmissionRequestIndexAction) loadParams() {
	action.ValidateCursorAsDefault()
	action.ExchangeFilter = action.GetString("exchange")
	action.StateFilter = action.GetInt32("state")
	action.ReferenceFilter = action.GetString("reference")
	action.AssetFilter = action.GetString("asset")
	action.PagingParams = action.GetPageQuery()
	action.Page.Filters = map[string]string{
		"state":     strconv.FormatInt(int64(action.StateFilter), 10),
		"asset":     action.AssetFilter,
		"exchange":  action.ExchangeFilter,
		"reference": action.ReferenceFilter,
	}

}

func (action *CoinsEmissionRequestIndexAction) loadRecords() {
	requests := action.HistoryQ().CoinsEmissionRequests()

	if action.ExchangeFilter != "" {
		requests.ForExchange(action.ExchangeFilter)
	}

	if action.ReferenceFilter != "" {
		requests.ForReference(action.ReferenceFilter)
	}

	if action.AssetFilter != "" {
		requests.ForAsset(action.AssetFilter)
	}

	switch action.StateFilter {
	case REQUEST_APPROVED:
		state := true
		requests.ForState(&state)
	case REQUEST_REJECTED:
		state := false
		requests.ForState(&state)
	case REQUEST_PENDING:
		requests.ForState(nil)
	default: //no-op
	}

	err := requests.Page(action.PagingParams).Select(&action.Records)

	if err != nil {
		action.Log.WithError(err).Error("failed to get emission requests")
		action.Err = &problem.ServerError
		return
	}
}

func (action *CoinsEmissionRequestIndexAction) loadExchangeNames() {
	exchanges := map[string]string{}
	for i := range action.Records {
		// if exchange name is empty set it's ID
		exchanges[action.Records[i].Issuer] = action.Records[i].Issuer
	}

	// load all exchanges
	for exchangeKey := range exchanges {
		exchangeName, err := action.CoreQ().ExchangeName(exchangeKey)
		if err != nil {
			action.Log.WithError(err).Error("Failed to get exchange name")
			action.Err = &problem.ServerError
			return
		}

		if exchangeName == nil {
			continue
		}

		exchanges[exchangeKey] = *exchangeName
	}

	// populate names
	for i := range action.Records {
		action.Records[i].ExchangeName = exchanges[action.Records[i].Issuer]
	}
}

func (action *CoinsEmissionRequestIndexAction) loadPage() {
	for _, record := range action.Records {
		var request resource.CoinsEmissionRequest
		request.Populate(&record)
		action.Page.Add(request)
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}

func (action *CoinsEmissionRequestIndexAction) checkAllowed() {
	action.IsAllowed(action.ExchangeFilter)
}

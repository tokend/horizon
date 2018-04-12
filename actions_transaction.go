package horizon

import (
	"net/http"

	"github.com/go-errors/errors"
	"gitlab.com/distributed_lab/logan"
	"gitlab.com/distributed_lab/txsub"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/render/sse"
	"gitlab.com/swarmfund/horizon/resource"
	txsubHelper "gitlab.com/swarmfund/horizon/txsub"
)

// This file contains the actions:
//
// TransactionIndexAction: pages of transactions
// TransactionShowAction: single transaction by sequence, by hash or id

// TransactionIndexAction renders a page of ledger resources, identified by
// a normal page query.
type TransactionIndexAction struct {
	Action
	LedgerFilter  int32
	AccountFilter string
	BalanceFilter string
	PagingParams  db2.PageQuery
	Records       []history.Transaction
	MetaLedger    history.Ledger
	Page          hal.Page
}

// JSON is a method for actions.JSON
func (action *TransactionIndexAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.checkAllowed,
		action.ValidateCursorWithinHistory,
		action.loadRecords,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

// SSE is a method for actions.SSE
func (action *TransactionIndexAction) SSE(stream sse.Stream) {
	action.Setup(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.checkAllowed,
		action.ValidateCursorWithinHistory,
	)
	action.Do(
		func() {
			// we will reuse this variable in sse, so re-initializing is required
			action.Records = []history.Transaction{}
		},
		action.loadRecords,
		func() {
			records := action.Records[:]

			for _, record := range records {
				var res resource.Transaction
				res.Populate(action.Ctx, record)
				stream.Send(sse.Event{
					ID:   res.PagingToken(),
					Data: res,
				})
				action.PagingParams.Cursor = res.PagingToken()
			}
		},
	)
}

const (
	maxTxPagSize uint64 = 1000
)

func (action *TransactionIndexAction) loadParams() {
	action.ValidateCursorAsDefault()
	action.AccountFilter = action.GetString("account_id")
	action.LedgerFilter = action.GetInt32("ledger_id")

	action.PagingParams = action.GetPageQuery()

	if action.PagingParams.Limit > maxTxPagSize {
		action.PagingParams.Limit = maxTxPagSize
	}
}

func (action *TransactionIndexAction) loadRecords() {
	q := action.HistoryQ()
	txs := q.Transactions()

	switch {
	case action.AccountFilter != "":
		txs.ForAccount(action.AccountFilter)
	case action.LedgerFilter > 0:
		txs.ForLedger(action.LedgerFilter)
	}

	// memorize ledger sequence before select to prevent data race
	latestLedger := int32(action.App.historyLatestLedgerGauge.Value())

	err := txs.Page(action.PagingParams).Select(&action.Records)
	if err != nil {
		action.Log.WithError(err).Error("failed to get transactions")
		action.Err = &problem.ServerError
		return
	}

	if uint64(len(action.Records)) == action.PagingParams.Limit {
		// we fetched full page, probably there is something ahead
		latestLedger = action.Records[len(action.Records)-1].LedgerSequence
	}

	// load ledger close time
	if err := action.HistoryQ().LedgerBySequence(&action.MetaLedger, latestLedger); err != nil {
		action.Log.WithError(err).Error("failed to get ledger")
		action.Err = &problem.ServerError
		return
	}
}

func (action *TransactionIndexAction) loadPage() {
	for _, record := range action.Records {
		var res resource.Transaction
		res.Populate(action.Ctx, record)
		action.Page.Add(res)
	}

	action.Page.Embedded.Meta = &hal.PageMeta{
		LatestLedger: &hal.LatestLedgerMeta{
			Sequence: action.MetaLedger.Sequence,
			ClosedAt: action.MetaLedger.ClosedAt,
		},
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}

func (action *TransactionIndexAction) checkAllowed() {
	action.IsAllowed(action.AccountFilter)
}

// TransactionShowAction renders a ledger found by its sequence number.
type TransactionShowAction struct {
	Action
	HashOrID string
	Record   history.Transaction
	Resource resource.Transaction
}

func (action *TransactionShowAction) loadParams() {
	action.HashOrID = action.GetString("id")
}

func (action *TransactionShowAction) loadRecord() {
	action.Err = action.HistoryQ().TransactionByHashOrID(&action.Record, action.HashOrID)
}

func (action *TransactionShowAction) loadResource() {
	action.Resource.Populate(action.Ctx, action.Record)
}

// JSON is a method for actions.JSON
func (action *TransactionShowAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.checkAllowed,
		action.loadRecord,
		action.loadResource,
		func() { hal.Render(action.W, action.Resource) },
	)
}

func (action *TransactionShowAction) checkAllowed() {
	action.IsAllowed("")
}

// TransactionCreateAction submits a transaction to the stellar-core network
// on behalf of the requesting client.
type TransactionCreateAction struct {
	Action
	TX       string
	Result   txsub.Result
	Resource resource.TransactionSuccess

	TFAFailed bool
}

// JSON format action handler
func (action *TransactionCreateAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadTX,
		action.checkAllowed,
		action.loadResult,
		action.loadResource,
		func() {
			if !action.TFAFailed {
				hal.Render(action.W, action.Resource)
			}
		})
}

func (action *TransactionCreateAction) loadTX() {
	action.TX = action.GetNonEmptyString("tx")
}

func (action *TransactionCreateAction) checkAllowed() {
	if action.App.config.DisableAPISubmit {
		return
	}
	action.isAllowed("")
}

func (action *TransactionCreateAction) loadResult() {
	if action.TFAFailed {
		return
	}

	envelopeInfo, err := txsubHelper.ExtractEnvelopeInfo(action.TX, action.App.CoreInfo.NetworkPassphrase)
	if err != nil {
		action.Err = &problem.P{
			Type:   "transaction_malformed",
			Title:  "Transaction Malformed",
			Status: http.StatusBadRequest,
			Detail: "Horizon could not decode the transaction envelope in this " +
				"request. A transaction should be an XDR TransactionEnvelope struct " +
				"encoded using base64.  The envelope read from this request is " +
				"echoed in the `extras.envelope_xdr` field of this response for your " +
				"convenience.",
			Extras: map[string]interface{}{
				"envelope_xdr": action.TX,
			},
		}
		return
	}

	action.Result = action.App.submitter.Submit(action.Ctx, envelopeInfo)
	if action.Result.HasInternalError() {
		action.Log.WithError(action.Result.Err).Error("Failed to submit tx")
		action.Err = &problem.ServerError
		return
	}

	if action.Result.Err == nil {
		action.Resource.Populate(action.Ctx, action.Result)
		return
	}
}

func (action *TransactionCreateAction) loadResource() {
	if action.TFAFailed {
		return
	}

	p, err := txResultToProblem(&action.Result)
	if err != nil {
		action.Log.WithError(err).Error("failed to craft problem")
		action.Err = &problem.ServerError
		return
	}

	if p != nil {
		action.Err = p
		return
	}
}

func txResultToProblem(result *txsub.Result) (*problem.P, error) {
	if result.Err == nil {
		return nil, nil
	}

	txSubError, ok := result.Err.(txsub.Error)
	if !ok {
		return nil, errors.New("Unexpected error type")
	}

	switch txSubError.Type() {
	case txsub.Timeout:
		return &problem.Timeout, nil
	case txsub.RejectedTx:
		var parsedResult xdr.TransactionResult
		err := xdr.SafeUnmarshalBase64(txSubError.ResultXDR(), &parsedResult)
		if err != nil {
			return nil, logan.Wrap(err, "Failed to get parse tx result")
		}

		var parsedEnvelope xdr.TransactionEnvelope
		err = xdr.SafeUnmarshalBase64(result.EnvelopeXDR, &parsedEnvelope)
		if err != nil {
			return nil, logan.Wrap(err, "Failed to unmarshal tx envelope")
		}

		resultCodes, err := resource.NewTransactionResultCodes(parsedResult)
		if err != nil {
			return nil, logan.Wrap(err, "Failed to create transaction result codes")
		}

		return &problem.P{
			Type:   "transaction_failed",
			Title:  "Transaction Failed",
			Status: http.StatusBadRequest,
			Detail: "The transaction failed when submitted to the stellar network. " +
				"The `extras.result_codes` field on this response contains further " +
				"details.  Descriptions of each code can be found at: " +
				"https://www.stellar.org/developers/learn/concepts/list-of-operations.html",
			Extras: map[string]interface{}{
				"envelope_xdr":    result.EnvelopeXDR,
				"result_xdr":      txSubError.ResultXDR(),
				"result_codes":    resultCodes,
				"parsed_result":   &parsedResult,
				"parsed_envelope": &parsedEnvelope,
			},
		}, nil
	default:
		return nil, errors.Errorf("Unexpected error type: %d", txSubError.Type())
	}
}

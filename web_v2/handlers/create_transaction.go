package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/google/jsonapi"

	"gitlab.com/tokend/regources/generated"

	"gitlab.com/tokend/horizon/web_v2/resources"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/tokend/horizon/db2/history2"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/txsub/v2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
)

//CreateTransaction submits transaction to the core
func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCreateTransactionRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	coreRepo := ctx.CoreRepo(r)
	historyRepo := ctx.HistoryRepo(r)
	results := txsub.ResultsProvider{
		Core:    core2.NewTransactionQ(coreRepo),
		History: history2.NewTransactionsQ(historyRepo),
	}
	handler := createTransactionHandler{
		Results:   results,
		Log:       ctx.Log(r),
		Submitter: ctx.Submitter(r),
		LedgerQ:   *history2.NewLedgerQ(historyRepo),
	}

	handler.Log = handler.Log.WithFields(logan.F{
		"tx_hash": request.Env.ContentHash,
	})

	result, err := handler.createTx(r.Context(), request)
	if errObj, ok := err.(*jsonapi.ErrorObject); ok {
		ape.RenderErr(w, errObj)
		return
	}
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to create transaction ", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if result == nil {
		ctx.Log(r).Error("got empty result", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type createTransactionHandler struct {
	Results   txsub.ResultsProvider
	LedgerQ   history2.LedgerQ
	Submitter *txsub.System
	Log       *logan.Entry
}

func (h *createTransactionHandler) createTx(context context.Context, request *requests.CreateTransaction) (*regources.TransactionResponse, error) {
	res, err := h.Submitter.Submit(context, *request.Env)
	if txsub.IsInternalError(err) {
		return nil, resources.NewTxFailure(*request.Env, err.(txsub.Error))
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to handle create transaction request")
	}
	if res == nil {
		return nil, errors.New("failed to submit transaction")
	}

	if !request.WaitForIngest {
		return h.prepareFromResult(res), nil
	}

	return h.waitIngest(context, res)
}

func (h *createTransactionHandler) tryGetFromHistory(hash string) (*history2.Transaction, bool) {
	tx, err := h.Results.History.GetByHash(hash)
	if err != nil {
		h.Log.
			WithError(err).
			Error("failed to get tx by hash")
		return nil, false
	}

	return tx, tx != nil
}

func (h *createTransactionHandler) waitIngest(context context.Context, result *txsub.Result) (*regources.TransactionResponse, error) {
	tx, ok := h.tryGetFromHistory(result.Hash)
	if ok {
		return h.prepareFromHistory(tx), nil
	}
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-context.Done():
			return nil, nil
		case <-ticker.C:
			tx, ok = h.tryGetFromHistory(result.Hash)
			if ok {
				return h.prepareFromHistory(tx), nil
			}
		}
	}
}

func (h *createTransactionHandler) prepareFromHistory(transaction *history2.Transaction) *regources.TransactionResponse {
	response := &regources.TransactionResponse{}
	response.Data = resources.NewTransaction(*transaction)
	meta := h.getTransactionMeta()
	if meta != nil {
		response.Meta = *meta
	}
	return response
}

func (h *createTransactionHandler) prepareFromResult(result *txsub.Result) *regources.TransactionResponse {
	response := &regources.TransactionResponse{}
	data := regources.Transaction{
		Key: resources.NewTxKey(result.TransactionID),
		Attributes: regources.TransactionAttributes{
			EnvelopeXdr:    result.EnvelopeXDR,
			Hash:           result.Hash,
			LedgerSequence: result.LedgerSequence,
			ResultMetaXdr:  result.ResultMetaXDR,
			ResultXdr:      result.ResultXDR,
		},
	}
	response.Data = data
	meta := h.getTransactionMeta()
	if meta != nil {
		response.Meta = *meta
	}
	return response
}

func (h *createTransactionHandler) getLatestLedger() (*history2.Ledger, error) {
	sequence, err := h.LedgerQ.GetLatestLedgerSeq()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load latest ledger sequence")
	}

	ledger, err := h.LedgerQ.GetBySequence(sequence)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load ledger by sequence")
	}

	return ledger, nil
}

func (h *createTransactionHandler) getTransactionMeta() *regources.TransactionResponseMeta {
	ledger, err := h.getLatestLedger()
	if err != nil {
		h.Log.WithError(err).Error("Failed to get latest ledger")
		return nil
	}

	return &regources.TransactionResponseMeta{
		LatestLedgerCloseTime: ledger.ClosedAt,
		LatestLedgerSequence:  ledger.Sequence,
	}
}

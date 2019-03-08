package handlers

import (
	"context"
	"net/http"
	"time"

	"gitlab.com/tokend/regources/v2"

	"gitlab.com/tokend/horizon/web_v2/resources"

	"gitlab.com/tokend/horizon/ingest2/storage"

	"github.com/lib/pq"

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
	coreRepo := ctx.CoreRepo(r)
	historyRepo := ctx.HistoryRepo(r)
	results := txsub.ResultsProvider{
		Core:    core2.NewTransactionQ(coreRepo),
		History: history2.NewTransactionQ(historyRepo),
	}
	config := ctx.Config(r)
	handler := createTransactionHandler{
		Results:       results,
		Log:           ctx.Log(r),
		Txsub:         ctx.Submitter(r),
		WaitForIngest: config.Ingest && config.WaitForIngest,
		Listener: pq.NewListener(config.DatabaseURL, 3*time.Second, 8*time.Second, func(event pq.ListenerEventType, err error) {
		}),
	}

	request, err := requests.NewCreateTransactionRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	handler.Log = handler.Log.WithFields(logan.F{
		"tx_hash": request.Env.ContentHash,
	})

	result, err := handler.createTx(r.Context(), request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to create transaction ", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if result == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if result.TxSubmitFailure != nil {
		w.WriteHeader(result.TxSubmitFailure.Status)
	}
	ape.Render(w, result)
}

type createTransactionHandler struct {
	Results       txsub.ResultsProvider
	Txsub         *txsub.System
	Listener      *pq.Listener
	WaitForIngest bool
	Log           *logan.Entry
}

func (h *createTransactionHandler) createTx(context context.Context, request *requests.CreateTransaction) (*regources.TxSubmitResponse, error) {
	res, err := h.Txsub.Submit(context, *request.Env)
	if txsub.IsInternalError(err) {
		return resources.NewTxFailure(*request.Env, err.(txsub.Error)), nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to handle create transaction request")
	}
	if res == nil {
		return nil, errors.New("failed to submit transaction")
	}
	if !h.WaitForIngest {
		return resources.NewTxSuccess(res), nil
	}

	err = h.Listener.Listen(storage.ChanSubmitter)
	if err != nil {
		h.Log.WithError(errors.Wrap(err, "failed to listen channel", logan.F{
			"channel": storage.ChanSubmitter,
		})).Error("Got error while waiting for tx ingest")
		return resources.NewTxSuccess(res), nil
	}

waitForIngest:
	for {
		select {
		case <-context.Done():
			break waitForIngest
		case <-h.Listener.Notify:
			tx, err := h.Results.History.GetByHash(res.Hash)
			if err != nil {
				h.Log.
					WithError(errors.Wrap(err, "failed to get tx by hash ")).
					Error("Got error while waiting for tx ingestion ")
				return resources.NewTxSuccess(res), nil
			}
			if tx != nil {
				break waitForIngest
			}
		}
	}

	return resources.NewTxSuccess(res), nil

}

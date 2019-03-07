package handlers

import (
	"context"
	"net/http"

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
	handler := createTransactionHandler{
		Results: results,
		Log:     ctx.Log(r),
		Txsub:   ctx.Submitter(r),
	}

	request, err := requests.NewCreateTransactionRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.createTx(r.Context(), request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to create transaction", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if result == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, result)
}

type createTransactionHandler struct {
	Results txsub.ResultsProvider
	Txsub   *txsub.System
	Log     *logan.Entry
}

func (h *createTransactionHandler) createTx(ctx context.Context, request *requests.CreateTransaction) (*txsub.Result, error) {
	res, err := h.Txsub.Submit(ctx, *request.Env)
	if err != nil {
		return nil, errors.Wrap(err, "failed to handle create transaction request")
	}

	return res, nil
}

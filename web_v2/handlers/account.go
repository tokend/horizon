package handlers

import (
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/horizon/web_v2/resource"
	"net/http"
)

func ShowAccount(w http.ResponseWriter, r *http.Request) {
	handler := ShowAccountHandler{}
	handler.Prepare(r, w)

	handler.accountId = chi.URLParam(r, "id")

	handler.Serve()
}

type ShowAccountHandler struct {
	Base

	accountId string
}

func (h *ShowAccountHandler) Serve() {
	coreAccount, err := h.CoreQ.
		Accounts().
		ByAddress(h.accountId)
	if err != nil {
		h.Logger.WithError(err).Error("Failed to get account record")
		h.RenderErr(problems.InternalError())
		return
	}
	account := resource.NewAccount(coreAccount)

	coreBalances, err := h.CoreQ.
		Balances().
		ByAddress(h.accountId).
		Select()
	if err != nil {
		h.Logger.WithError(err).Error("Failed to get balance records")
		h.RenderErr(problems.InternalError())
		return
	}
	balances := resource.NewBalanceCollection(coreBalances)

	account.Data.RelateBalances(*balances)
	if h.IsRequested("balances") {
		account.IncludeBalances(balances.Data)
	}

	h.Render(account)
}

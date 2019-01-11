package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/web_v2/resource"
)

func ShowAccountCollection(w http.ResponseWriter, r *http.Request) {
	handler := ShowAccountCollectionHandler{}
	handler.Prepare(r, w)
	handler.Serve()
}

type ShowAccountCollectionHandler struct {
	Base
}

func (h *ShowAccountCollectionHandler) Serve() {
	err := h.GetPageQuery()
	if err != nil {
		h.Logger.WithError(err).Error("Failed to get account records")
		h.RenderErr(problems.InternalError())
		return
	}

	var coreAccounts []core.Account
	err = h.CoreQ.
		Accounts().
		PageV2(h.PageQuery).
		Select(&coreAccounts)
	if err != nil {
		h.Logger.WithError(err).Error("Failed to get account records")
		h.RenderErr(problems.InternalError())
		return
	}
	accounts := resource.NewAccountCollection(coreAccounts)

	for i, coreAccount := range coreAccounts {
		coreBalances, err := h.CoreQ.
			Balances().
			ByAddress(coreAccount.AccountID).
			Select()
		if err != nil {
			h.Logger.WithError(err).Error("Failed to get balances for" + coreAccount.AccountID)
			h.RenderErr(problems.InternalError())
			return
		}

		balances := resource.NewBalanceCollection(coreBalances)
		accounts.Data[i].RelateBalances(*balances)

		if h.IsRequested("balances") {
			accounts.IncludeBalances(balances.Data)
		}
	}

	accounts.Links = h.GetLinksObject()
	h.Render(accounts)
}

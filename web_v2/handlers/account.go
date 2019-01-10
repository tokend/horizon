package handlers

import (
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/resource"
	"net/http"
)

func ShowAccount(w http.ResponseWriter, r *http.Request) {
	// TODO: isAllowed check
	accountId := chi.URLParam(r, "id")

	coreAccount, err := ctx.CoreQ(r).Accounts().ByAddress(accountId)
	if err != nil {
		ctx.Log(r).WithError(err).Error("Failed to get account record")
		RenderErr(w, r, problems.InternalError())
	}
	account := resource.NewAccount(coreAccount)

	coreBalances, err := ctx.CoreQ(r).Balances().ByAddress(accountId).Select()
	if err != nil {
		ctx.Log(r).WithError(err).Error("Failed to get account record")
		RenderErr(w, r, problems.InternalError())
	}
	balances := resource.NewBalanceCollection(coreBalances)
	account.Data.RelateBalances(*balances)
	if Includes("balances") {
		account.IncludeBalances(balances.Data)
	}

	err = Render(w, account)
	if err != nil {
		RenderErr(w, r, problems.InternalError())
	}
}

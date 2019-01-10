package handlers

import (
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/resource"
	"net/http"
)

func ShowAccountCollection(w http.ResponseWriter, r *http.Request) {
	var coreAccounts []core.Account

	err := ctx.CoreQ(r).Accounts().Select(&coreAccounts)
	if err != nil {
		ctx.Log(r).WithError(err).Error("Failed to get account records")
		RenderErr(w, r, problems.InternalError())
		return
	}

	accounts := resource.NewAccountCollection(coreAccounts)

	for i, coreAccount := range coreAccounts {
		coreBalances, err := ctx.CoreQ(r).Balances().ByAddress(coreAccount.AccountID).Select()
		if err != nil {
			ctx.Log(r).WithError(err).Error("Failed to get balances for" + coreAccount.AccountID)
			RenderErr(w, r, problems.InternalError())
			return
		}

		balances := resource.NewBalanceCollection(coreBalances)
		accounts.Data[i].RelateBalances(*balances)

		if Includes("balances") {
			accounts.IncludeBalances(balances.Data)
		}
	}

	err = Render(w, accounts)
	if err != nil {
		RenderErr(w, r, problems.InternalError())
	}
}

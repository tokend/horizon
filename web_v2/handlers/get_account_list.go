package handlers

import (
	"net/http"

	"gitlab.com/tokend/horizon/web_v2/ctx"

	"gitlab.com/tokend/horizon/web_v2/resources"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/tokend/horizon/web_v2/requests"
	regources "gitlab.com/tokend/regources/generated"

	"gitlab.com/tokend/horizon/db2/core2"
)

// GetAccountList - get list of accounts by address
func GetAccountList(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getAccountListHandler{
		AccountsQ: core2.NewAccountsQ(coreRepo),
	}

	request, err := requests.NewGetAccountList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	response, err := handler.GetAccountList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get account", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, response)
}

type getAccountListHandler struct {
	AccountsQ core2.AccountsQ
}

func (h *getAccountListHandler) GetAccountList(r *requests.GetAccountList) (*regources.AccountListResponse, error) {
	q := h.AccountsQ.Page(r.PageParams)

	if r.ShouldFilter(requests.FilterTypeAccountListAccount) {
		q = q.FilterByAddresses(r.Filters.Account)
	}
	if r.ShouldFilter(requests.FilterTypeAccountListRole) {
		q = q.FilterByRole(r.Filters.Role)
	}

	accounts, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get accounts")
	}

	response := regources.AccountListResponse{
		Data:  make([]regources.Account, 0, len(accounts)),
		Links: r.GetOffsetLinks(r.PageParams),
	}

	for _, account := range accounts {
		accountResource := resources.NewAccount(account, nil)
		accountResource.Relationships.Role = resources.NewAccountRoleKey(account.RoleID).AsRelation()
		response.Data = append(response.Data, accountResource)
	}

	return &response, nil
}

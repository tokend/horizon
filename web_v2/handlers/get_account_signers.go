package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"
	"gitlab.com/tokend/regources/v2"
)

// GetAccountSigners - processes request to get account signers
func GetAccountSigners(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getAccountSignersHandler{
		SignersQ: core2.NewSignerQ(coreRepo),
		AccountQ: core2.NewAccountsQ(coreRepo),
		Log:      ctx.Log(r),
	}

	request, err := requests.NewGetAccountSigners(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetAccountSigners(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get account", logan.F{
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

type getAccountSignersHandler struct {
	SignersQ core2.SignerQ
	AccountQ core2.AccountsQ
	Log      *logan.Entry
}

//GetAccountSigners - returns signers for account
func (h *getAccountSignersHandler) GetAccountSigners(request *requests.GetAccountSigners) ([]*regources.Signer, error) {
	account, err := h.AccountQ.GetByAddress(request.Address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load account", logan.F{
			"account_address": request.Address,
		})
	}

	if account == nil {
		return nil, nil
	}

	signers, err := h.SignersQ.FilterByAccountAddress(request.Address).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load signers for account", logan.F{
			"account_address": request.Address,
		})
	}

	// TODO: KILL it
	if account.Thresholds[0] > 0 {
		signers = append(signers, core2.Signer{
			ID:        account.Address,
			AccountID: account.Address,
			Weight:    int(account.Thresholds[0]),
			Name:      "master",
		})
	}

	response := make([]*regources.Signer, 0, len(signers))
	for i := range signers {
		signer := resources.NewSigner(signers[i])
		if request.NeedRoles() {
			signer.Role = h.getRole(request)
		}
		response = append(response, signer)
	}

	return response, nil
}

func (h *getAccountSignersHandler) getRole(request *requests.GetAccountSigners) *regources.Role {
	result := regources.Role{
		ID: "mocked_role",
		Details: map[string]interface{}{
			"name": "Name of the Mocked Role",
		},
	}

	if !request.NeedRules() {
		return &result
	}

	result.Rules = []*regources.Rule{
		{
			ID:       "mocked_rule_id",
			Resource: "NOTE: format will be changed",
			Action:   "view",
			Details: map[string]interface{}{
				"name": "Name of the mocked Rule",
			},
		},
	}

	return &result
}

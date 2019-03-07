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

// GetPublicKeyAccountList - processes request to get the list of all accounts that has provided public key as it's signer
func GetPublicKey(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getPublicKeyAccountListHandler{
		SignersQ: core2.NewSignerQ(coreRepo),
		Log:       ctx.Log(r),
	}

	request, err := requests.NewGetPublicKey(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !isAllowed(r, w, request.ID) {
		return
	}

	result, err := handler.GetPublicKeyAccountList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get account list for public key", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getPublicKeyAccountListHandler struct {
	SignersQ core2.SignerQ
	Log       *logan.Entry
}

// GetAccountRuleList returns the list of accountRules with related resources
func (h *getPublicKeyAccountListHandler) GetPublicKeyAccountList(request *requests.GetPublicKey) (*regources.PublicKeyResponse, error) {
	coreSigners, err := h.SignersQ.FilterByPublicKey(request.ID).Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get signers list")
	}

	accounts := make([]regources.Key, 0, len(coreSigners))
	for _, coreSigner := range coreSigners {
		accounts = append(accounts, resources.NewAccountKey(coreSigner.AccountID))
	}

	publicKeyResource := resources.NewPublicKeyResource(request.ID)
	publicKeyResource.Relationships.Accounts = &regources.RelationCollection{
		Data: accounts,
	}

	response := &regources.PublicKeyResponse{
		Data: publicKeyResource,
	}

	return response, nil
}

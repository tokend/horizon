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

// GetSignerRuleList - processes request to get the list of signerRules
func GetSignerRuleList(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getSignerRuleListHandler{
		SignerRulesQ: core2.NewSignerRuleQ(coreRepo),
		Log:          ctx.Log(r),
	}

	request, err := requests.NewGetSignerRuleList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetSignerRuleList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get signerRule list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getSignerRuleListHandler struct {
	SignerRulesQ core2.SignerRuleQ
	Log          *logan.Entry
}

// GetSignerRuleList returns the list of signerRules with related resources
func (h *getSignerRuleListHandler) GetSignerRuleList(request *requests.GetSignerRuleList) (*regources.SignerRulesResponse, error) {
	signerRules, err := h.SignerRulesQ.Page(*request.PageParams).Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get signer rule list")
	}

	response := &regources.SignerRulesResponse{
		Data:  make([]regources.SignerRule, 0, len(signerRules)),
		Links: request.GetOffsetLinks(*request.PageParams),
	}

	for _, signerRule := range signerRules {
		signerRuleResponse := resources.NewSignerRule(signerRule)
		response.Data = append(response.Data, signerRuleResponse)
	}

	return response, nil
}

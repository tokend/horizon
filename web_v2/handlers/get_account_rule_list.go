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
	"gitlab.com/tokend/regources/rgenerated"
)

// GetAccountRuleList - processes request to get the list of accountRules
func GetAccountRuleList(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getAccountRuleListHandler{
		AccountRulesQ: core2.NewAccountRuleQ(coreRepo),
		Log:           ctx.Log(r),
	}

	request, err := requests.NewGetAccountRuleList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetAccountRuleList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get accountRule list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getAccountRuleListHandler struct {
	AccountRulesQ core2.AccountRuleQ
	Log           *logan.Entry
}

// GetAccountRuleList returns the list of accountRules with related resources
func (h *getAccountRuleListHandler) GetAccountRuleList(request *requests.GetAccountRuleList) (*rgenerated.AccountRulesResponse, error) {
	accountRules, err := h.AccountRulesQ.Page(*request.PageParams).Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get account rule list")
	}

	response := &rgenerated.AccountRulesResponse{
		Data:  make([]rgenerated.AccountRule, 0, len(accountRules)),
		Links: request.GetOffsetLinks(*request.PageParams),
	}

	for _, accountRule := range accountRules {
		accountRuleResponse := resources.NewAccountRule(accountRule)
		response.Data = append(response.Data, accountRuleResponse)
	}

	return response, nil
}

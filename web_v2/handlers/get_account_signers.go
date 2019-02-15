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
	SignersQ    core2.SignerQ
	AccountQ    core2.AccountsQ
	SignerRoleQ core2.SignerRoleQ
	Log         *logan.Entry
}

//GetAccountSigners - returns signers for account
func (h *getAccountSignersHandler) GetAccountSigners(request *requests.GetAccountSigners) (*regources.SignersResponse, error) {
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

	response := regources.SignersResponse{
		Data: make([]regources.Signer, 0, len(signers)),
	}
	for _, signerRaw := range signers {
		signer := resources.NewSigner(signerRaw)
		if request.ShouldInclude(requests.IncludeTypeSignerRoles) {
			signer.Relationships.Role, err = h.getRole(request, &response.Included, signerRaw)

			if err != nil {
				return nil, errors.Wrap(err, "failed to get role")
			}
		}
		response.Data = append(response.Data, signer)
	}

	return &response, nil
}

func (h *getAccountSignersHandler) getRole(request *requests.GetAccountSigners,
	includes *regources.Included, signer core2.Signer,
) (*regources.Relation, error) {
	if !request.ShouldInclude(requests.IncludeTypeSignerRoles) {
		role := resources.NewSignerRoleKey(signer.RoleID)
		return role.AsRelation(), nil
	}

	roleRaw, err := h.SignerRoleQ.GetByID(signer.RoleID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get signer role")
	}

	if roleRaw == nil {
		return nil, errors.New("signer role not found")
	}

	role := resources.NewSignerRole(*roleRaw)

	if request.ShouldInclude(requests.IncludeTypeSignerRolesRules) {
		rules := []regources.SignerRule{
			{
				Key: regources.Key{
					ID:   "mocked_rule_id",
					Type: regources.TypeSignerRules,
				},
				Attributes: regources.SignerRuleAttr{
					Resource: "NOTE: format will be changed",
					Action:   "view",
					Details:  []byte{},
				},
			},
		}

		role.Relationships.Rules = &regources.RelationCollection{}
		for i := range rules {
			role.Relationships.Rules.Data = append(role.Relationships.Rules.Data, rules[i].GetKey())
			includes.Add(&rules[i])
		}
	}

	includes.Add(&role)

	return role.AsRelation(), nil
}

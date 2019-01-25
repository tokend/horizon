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

// GetAssetList - processes request to get the list of assets
func GetAssetList(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getAssetListHandler{
		AccountsQ: core2.NewAccountsQ(coreRepo),
		AssetsQ:   core2.NewAssetsQ(coreRepo),
		Log:       ctx.Log(r),
	}

	request, err := requests.NewGetAssetList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetAssetList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get asset list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getAssetListHandler struct {
	AssetsQ   core2.AssetsQ
	AccountsQ core2.AccountsQ
	Log       *logan.Entry
}

// GetAssetList returns the list of assets with related resources
func (h *getAssetListHandler) GetAssetList(request *requests.GetAssetList) (*regources.AssetsResponse, error) {
	q := h.AssetsQ.Page(*request.PageParams)
	if request.ShouldFilter(requests.FilterTypeAssetListOwner) {
		q = q.FilterByOwner(request.Filters.Owner)
	}
	if request.ShouldFilter(requests.FilterTypeAssetListPolicy) {
		q = q.FilterByPolicy(request.Filters.Policy)
	}
	assets, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get asset list")
	}

	response := &regources.AssetsResponse{
		Data:  make([]regources.Asset, 0, len(assets)),
		Links: request.GetOffsetLinks(*request.PageParams),
	}

	for i := range assets {
		asset := resources.NewAsset(assets[i])
		owner := resources.NewAccountKey(assets[i].Owner)
		asset.Relationships.Owner = owner.AsRelation()

		if request.ShouldInclude(requests.IncludeTypeAssetListOwners) {
			response.Included.Add(&owner)
		}

		response.Data = append(response.Data, asset)
	}

	return response, nil
}

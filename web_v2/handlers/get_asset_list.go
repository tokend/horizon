package handlers

import (
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"
	"gitlab.com/tokend/regources/v2"
	"net/http"
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

	pageParams, err := request.GetOffsetBasedPageParams()
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetAssetList(request, pageParams.Limit(), pageParams.Offset())
	if result == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get asset list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	links := pageParams.Links(request.URL())

	ape.RenderPage(w, result, links)
}

type getAssetListHandler struct {
	AssetsQ   core2.AssetsQ
	AccountsQ core2.AccountsQ
	Log       *logan.Entry
}

// GetAssetList returns the list of assets with related resources
func (h *getAssetListHandler) GetAssetList(request *requests.GetAssetList, limit, offset uint64) ([]*regources.Asset, error) {
	q := h.AssetsQ.Page(limit, offset)
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

	response := make([]*regources.Asset, 0, len(assets))
	for i := range assets {
		asset := resources.NewAsset(&assets[i])
		if request.ShouldInclude(requests.IncludeTypeAssetListOwners) {
			asset.Owner = &regources.Account{
				ID: assets[i].Owner,
			}
		}
		response = append(response, asset)
	}

	return response, nil
}

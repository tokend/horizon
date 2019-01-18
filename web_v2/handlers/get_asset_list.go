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

// GetAsset - processes request to get asset list
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

	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get asset list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	result, err := handler.GetAssetList(request)
	if result == nil {
		ape.RenderErr(w, problems.NotFound())
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
func (h *getAssetListHandler) GetAssetList(request *requests.GetAssetList) ([]*regources.Asset, error) {
	assets, err := h.AssetsQ.Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get asset by code")
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

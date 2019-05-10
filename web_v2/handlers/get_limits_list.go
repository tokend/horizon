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
	regources "gitlab.com/tokend/regources/generated"
)

// GetLimitsList - processes request to get the list of fees
func GetLimitsList(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getLimitsListHandler{
		LimitsQ: core2.NewLimitsQ(coreRepo),
		Log:     ctx.Log(r),
	}

	request, err := requests.NewGetLimitsList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetLimitsList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get fee list ", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getLimitsListHandler struct {
	LimitsQ  core2.LimitsQ
	AccountQ core2.AccountsQ
	AssetsQ  core2.AssetsQ
	Log      *logan.Entry
}

// GetLimitsList returns the list of fees with related resources
func (h *getLimitsListHandler) GetLimitsList(request *requests.GetLimitsList) (*regources.LimitssResponse, error) {
	q := h.LimitsQ.Page(*request.PageParams)
	if request.ShouldFilter(requests.FilterTypeLimitsListAccount) {
		q = q.FilterByAccount(request.Filters.Account)
	}
	if request.ShouldFilter(requests.FilterTypeLimitsListAccountRole) {
		q = q.FilterByAccountRole(request.Filters.AccountRole)
	}

	if request.ShouldFilter(requests.FilterTypeLimitsListAsset) {
		q = q.FilterByAsset(request.Filters.Asset)
	}

	if request.ShouldFilter(requests.FilterTypeLimitsListStatsOpType) {
		q = q.FilterByStatsOpType(request.Filters.StatsOpType)
	}

	limits, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get fee list")
	}

	response := &regources.LimitssResponse{
		Data:  make([]regources.Limits, 0, len(limits)),
		Links: request.GetOffsetLinks(*request.PageParams),
	}

	assets := make([]string, 0, len(limits))
	for i := range limits {
		limit := resources.NewLimits(limits[i])

		limit.Relationships.Asset = resources.NewAssetKey(limits[i].AssetCode).AsRelation()
		if limits[i].AccountId != nil {
			limit.Relationships.Account = resources.NewAccountKey(*limits[i].AccountId).AsRelation()
		}
		if limits[i].AccountType != nil {
			limit.Relationships.AccountRole = resources.NewAccountRoleKey(*limits[i].AccountType).AsRelation()
		}
		assets = append(assets, limits[i].AssetCode)
		response.Data = append(response.Data, limit)
	}

	if request.ShouldInclude(requests.IncludeTypeLimitsListAsset) {
		assetRecords, err := h.AssetsQ.FilterByCodes(assets).Select()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get assets")
		}
		if assetRecords == nil {
			return nil, errors.New("Assets not found")
		}

		for _, v := range assetRecords {
			asset := resources.NewAsset(v)
			response.Included.Add(&asset)
		}
	}

	return response, nil
}

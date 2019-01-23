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

// GetOfferList - processes request to get the list of offers
func GetOfferList(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getOfferListHandler{
		OffersQ:   core2.NewOffersQ(coreRepo),
		AssetsQ:   core2.NewAssetsQ(coreRepo),
		AccountsQ: core2.NewAccountsQ(coreRepo),
		Log:       ctx.Log(r),
	}

	request, err := requests.NewGetOfferList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetOfferList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get offer list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getOfferListHandler struct {
	OffersQ   core2.OffersQ
	AssetsQ   core2.AssetsQ
	AccountsQ core2.AccountsQ
	Log       *logan.Entry
}

// GetOfferList returns offer with related resources
func (h *getOfferListHandler) GetOfferList(request *requests.GetOfferList) (*regources.OffersResponse, error) {
	q := h.OffersQ.Page(request.PageParams.Limit(), request.PageParams.Offset())

	if request.ShouldFilter(requests.FilterTypeOfferListBaseBalance) {
		q = q.FilterByBaseBalanceID(request.Filters.BaseBalance)
	}

	if request.ShouldFilter(requests.FilterTypeOfferListQuoteBalance) {
		q = q.FilterByQuoteBalanceID(request.Filters.QuoteBalance)
	}

	if request.ShouldFilter(requests.FilterTypeOfferListBaseAsset) {
		q = q.FilterByBaseAssetCode(request.Filters.BaseAsset)
	}

	if request.ShouldFilter(requests.FilterTypeOfferListQuoteAsset) {
		q = q.FilterByQuoteAssetCode(request.Filters.QuoteAsset)
	}

	if request.ShouldFilter(requests.FilterTypeOfferListOwner) {
		q = q.FilterByOwnerID(request.Filters.Owner)
	}

	if request.ShouldFilter(requests.FilterTypeOfferListOrderBook) {
		q = q.FilterByOrderBookID(request.Filters.OrderBook)
	}

	if request.ShouldFilter(requests.FilterTypeOfferListIsBuy) {
		q = q.FilterByIsBuy(request.Filters.IsBuy)
	}

	if request.ShouldInclude(requests.IncludeTypeOfferListBaseAssets) {
		q = q.WithBaseAsset()
	}

	if request.ShouldInclude(requests.IncludeTypeOfferListQuoteAssets) {
		q = q.WithQuoteAsset()
	}

	coreOffers,err  := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get offer list")
	}

	response := &regources.OffersResponse{
		Data:  make([]regources.Offer, 0, len(coreOffers)),
		Links: request.PageParams.Links(request.URL()),
	}

	for _, coreOffer := range coreOffers {
		offer := resources.NewOffer(coreOffer)

		ownerKey := resources.NewAccountKey(coreOffer.OwnerID)
		baseAssetKey := resources.NewAssetKey(coreOffer.BaseAssetCode)
		quoteAssetKey := resources.NewAssetKey(coreOffer.QuoteAssetCode)
		baseBalanceKey := resources.NewBalanceKey(coreOffer.BaseBalanceID)
		quoteBalanceKey := resources.NewBalanceKey(coreOffer.QuoteBalanceID)

		offer.Relationships.Owner = ownerKey.AsRelation()
		offer.Relationships.BaseAsset = baseAssetKey.AsRelation()
		offer.Relationships.QuoteAsset = quoteAssetKey.AsRelation()
		offer.Relationships.BaseBalance = baseBalanceKey.AsRelation()
		offer.Relationships.QuoteBalance = quoteBalanceKey.AsRelation()

		response.Data = append(response.Data, offer)

		if request.ShouldInclude(requests.IncludeTypeOfferListBaseAssets) {
			coreBaseAsset := coreOffer.BaseAsset
			baseAsset := resources.NewAsset(*coreBaseAsset)

			response.Included.Add(&baseAsset)
		}

		if request.ShouldInclude(requests.IncludeTypeOfferListQuoteAssets) {
			coreQuoteAsset := coreOffer.QuoteAsset
			quoteAsset := resources.NewAsset(*coreQuoteAsset)

			response.Included.Add(&quoteAsset)
		}
	}

	return response, nil
}

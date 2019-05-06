package handlers

import (
	"fmt"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"
	"gitlab.com/tokend/regources/generated"
	"net/http"
)

// GetConvertedBalances - processes request to get converted balances and their details by accountID and asset code
func GetConvertedBalances(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)

	converter, err := newBalanceStateConverterForHandler(coreRepo)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed failed to create balance state converted")
		ape.Render(w, problems.InternalError())
		return
	}

	handler := getConvertedBalancesHandler{
		balanceStateConverter: converter,
		AssetsQ:               core2.NewAssetsQ(coreRepo),
		AccountsQ:             core2.NewAccountsQ(coreRepo),
		BalancesQ:             core2.NewBalancesQ(coreRepo),
		Log:                   ctx.Log(r),
	}

	request, err := requests.NewGetConvertedBalances(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetConvertedBalances(request)
	if err != nil {
		ctx.Log(r).WithError(err).WithField("request", request).Error("failed to get converted balances")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if result == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, result)
}

type getConvertedBalancesHandler struct {
	AssetsQ               core2.AssetsQ
	AccountsQ             core2.AccountsQ
	BalancesQ             core2.BalancesQ
	Log                   *logan.Entry
	balanceStateConverter *balanceStateConverter
}

// GetConvertedBalances - returns converted balances collection with related resources
func (h *getConvertedBalancesHandler) GetConvertedBalances(request *requests.GetConvertedBalances) (*regources.ConvertedBalancesCollectionResponse, error) {
	coreAccount, err := h.AccountsQ.GetByAddress(request.AccountAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account by address")
	}
	if coreAccount == nil {
		fmt.Println(request.AccountAddress)
		return nil, nil
	}

	coreAsset, err := h.AssetsQ.GetByCode(request.AssetCode)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get asset by code")
	}
	if coreAsset == nil {
		return nil, nil
	}

	coreBalances, err := h.BalancesQ.FilterByAccount(request.AccountAddress).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get balances by account address")
	}

	response := regources.ConvertedBalancesCollectionResponse{
		Data: resources.NewConvertedBalanceCollection(request.AssetCode),
	}
	response.Data.Relationships.States = regources.RelationCollection{
		Data: make([]regources.Key, 0, len(coreBalances)),
	}

	convertedStates := make([]regources.ConvertedBalanceState, 0, len(coreBalances))

	for _, coreBalance := range coreBalances {
		convertedState, err := h.balanceStateConverter.Convert(coreBalance, request.AssetCode)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get converted balance state")
		}
		convertedStates = append(convertedStates, *convertedState)
	}

	sortedConvertedStates := SortConvertedStates(convertedStates)

	for _, convertedState := range sortedConvertedStates {
		response.Data.Relationships.States.Data = append(
			response.Data.Relationships.States.Data,
			convertedState.Key,
		)

		if request.ShouldInclude(requests.IncludeTypeConvertedBalancesStates) {
			response.Included.Add(&convertedState)
		}
	}

	return &response, nil
}

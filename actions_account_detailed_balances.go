package horizon

import (
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/exchange"
	"github.com/go-errors/errors"
	"gitlab.com/swarmfund/go/amount"
)

type AccountDetailedBalancesAction struct {
	Action
	converter *exchange.Converter

	AccountID string
	ConvertToAsset string

	Balances []core.Balance
	Assets []core.Asset
	Sales []history.Sale

	AssetCodes []string

	Resource []resource.Balance
}

func (action *AccountDetailedBalancesAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadBalances,
		action.groupBalancesByAsset,
		action.loadAssets,
		action.createConverter,
		action.loadSales,
		action.loadResource,
		func () {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *AccountDetailedBalancesAction) loadParams() {
	action.AccountID = action.GetNonEmptyString("id")
}

func (action *AccountDetailedBalancesAction) checkAllowed() {
	action.IsAllowed(action.AccountID)
}

func (action *AccountDetailedBalancesAction) loadBalances() {
	var err error
	action.Balances, err = action.CoreQ().Balances().ByAddress(action.AccountID).Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to load balances")
		action.Err = &problem.ServerError
		return
	}
}

func (action *AccountDetailedBalancesAction) groupBalancesByAsset() {
	assetsMap := map[string]bool{}
	for _, balance := range action.Balances {
		if _, ok := assetsMap[balance.Asset]; ok {
			continue
		}

		assetsMap[balance.Asset] = true
		action.AssetCodes = append(action.AssetCodes, balance.Asset)
	}
}

func (action *AccountDetailedBalancesAction) loadAssets() {
	var err error
	action.Assets, err = action.CoreQ().Assets().ForCodes(action.AssetCodes).Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to load assets for codes")
		action.Err = &problem.ServerError
		return
	}
}

func (action *AccountDetailedBalancesAction) createConverter() {
	var err error
	action.converter, err = exchange.NewConverter(action.CoreQ())
	if err != nil {
		action.Log.WithError(err).Error("Failed to init converter")
		action.Err = &problem.ServerError
		return
	}
}

func (action *AccountDetailedBalancesAction) loadSales() {
	var err error
	action.Sales, err = selectSalesWithCurrentCap(action.HistoryQ().Sales().ForBaseAssets(action.AssetCodes...), action.converter)
	if err != nil {
		action.Log.WithError(err).Error("Failed to load sales")
		action.Err = &problem.ServerError
		return
	}
}

func (action *AccountDetailedBalancesAction) loadResource() {
	for _, record := range action.Balances {
		var r resource.Balance
		r.Populate(record)
		asset := findAssetByAssetCode(record.Asset, action.Assets)
		if asset == nil {
			action.Log.WithField("asset_code", record.Asset).Error("Failed to find asset for existing balance")
			action.Err = &problem.ServerError
			return
		}

		r.AssetDetails = asset
		r.AssetDetails.Sales = findAllSalesForAsset(asset.Code, action.Sales)

		var err error
		r.ConvertedBalance, err = convertAmount(record.Amount, r.Asset, action.ConvertToAsset, action.converter)
		if err != nil {
			action.Log.WithError(err).Error("Failed to convert balance")
			action.Err = &problem.ServerError
			return
		}

		r.ConvertedLocked, err = convertAmount(record.Locked, r.Asset, action.ConvertToAsset, converter)
		if err != nil {
			action.Log.WithError(err).Error("failed to convert locked amount")
			action.Err = &problem.ServerError
			return
		}


		action.Resource = append(action.Resource, r)
	}
}

func convertAmount(balance int64, fromAsset, toAsset string, converter *exchange.Converter) (string, error) {
	convertedAmount, err := converter.TryToConvertWithOneHop(balance, fromAsset, toAsset)
	if err != nil || convertedAmount == nil {
		if err == nil {
			err = errors.New("failed to find path to convert balance amount")
		}

		return "", err
	}

	return amount.String(*convertedAmount), nil
}

func findAssetByAssetCode(code string, assets []core.Asset) *resource.Asset {
	for i := range assets {
		if code != assets[i].Code {
			continue
		}

		var result resource.Asset
		result.Populate(&assets[i])
		return &result
	}

	return nil
}

func findAllSalesForAsset(code string, sales []history.Sale) []resource.Sale {
	var result []resource.Sale
	for i := range sales {
		if sales[i].BaseAsset != code {
			continue
		}

		var sale resource.Sale
		sale.Populate(&sales[i])
		result = append(result, sale)
	}

	return result
}



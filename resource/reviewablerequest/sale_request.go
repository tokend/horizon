package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources/reviewablerequest2"
)

func PopulateSaleCreationRequest(histRequest history.SaleRequest) (
	r *reviewablerequest2.SaleCreationRequest, err error,
) {
	r = &reviewablerequest2.SaleCreationRequest{}
	r.BaseAsset = histRequest.BaseAsset
	r.DefaultQuoteAsset = histRequest.DefaultQuoteAsset
	r.StartTime = histRequest.StartTime
	r.EndTime = histRequest.EndTime
	r.SoftCap = histRequest.SoftCap
	r.HardCap = histRequest.HardCap
	r.Details = histRequest.Details
	r.QuoteAssets = histRequest.QuoteAssets
	r.SaleType.Value = int32(histRequest.SaleType)
	r.SaleType.Name = histRequest.SaleType.ShortString()
	r.BaseAssetForHardCap = histRequest.BaseAssetForHardCap
	r.State.Value = int32(histRequest.State)
	r.State.Name = histRequest.State.ShortString()
	return
}

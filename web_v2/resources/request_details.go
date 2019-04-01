package resources

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/rgenerated"
)

// NewRequestDetails - returns new instance of reviewable request details
func NewRequestDetails(request history2.ReviewableRequest) rgenerated.Resource {
	switch request.RequestType {
	case xdr.ReviewableRequestTypeUpdateAsset:
		return newAssetUpdateRequest(request.ID, *request.Details.UpdateAsset)
	case xdr.ReviewableRequestTypeCreatePreIssuance:
		return newPreIssuanceRequest(request.ID, *request.Details.CreatePreIssuance)
	case xdr.ReviewableRequestTypeCreateIssuance:
		return newIssuanceRequest(request.ID, *request.Details.CreateIssuance)
	case xdr.ReviewableRequestTypeCreateWithdraw:
		return newWithdrawalRequest(request.ID, *request.Details.CreateWithdraw)
	case xdr.ReviewableRequestTypeCreateSale:
		return newSaleRequest(request.ID, *request.Details.CreateSale)
	case xdr.ReviewableRequestTypeUpdateLimits:
		return newLimitsUpdateRequest(request.ID, *request.Details.UpdateLimits)
	case xdr.ReviewableRequestTypeCreateAmlAlert:
		return newAmlAlertRequest(request.ID, *request.Details.CreateAmlAlert)
	case xdr.ReviewableRequestTypeChangeRole:
		return newChangeRoleRequest(request.ID, *request.Details.ChangeRole)
	case xdr.ReviewableRequestTypeUpdateSaleDetails:
		return newUpdateSaleDetailsRequest(request.ID, *request.Details.UpdateSaleDetails)
	case xdr.ReviewableRequestTypeCreateAsset:
		return newAssetCreateRequest(request.ID, *request.Details.CreateAsset)
	case xdr.ReviewableRequestTypeCreateAtomicSwapBid:
		return newAtomicSwapBidRequest(request.ID, *request.Details.CreateAtomicSwapBid)
	case xdr.ReviewableRequestTypeCreateAtomicSwap:
		return newAtomicSwapRequest(request.ID, *request.Details.CreateAtomicSwap)
	case xdr.ReviewableRequestTypeCreatePoll:
		return newCreatePollRequest(request.ID, *request.Details.CreatePoll)
	default:
		panic(errors.From(errors.New("unexpected operation type"), logan.F{
			"type": request.RequestType,
		}))
	}
	return nil
}

func newAmlAlertRequest(id int64, details history2.CreateAmlAlertRequest) *rgenerated.CreateAmlAlertRequest {
	return &rgenerated.CreateAmlAlertRequest{
		Key: rgenerated.NewKeyInt64(id, rgenerated.REQUEST_DETAILS_AML_ALERT),
		Attributes: rgenerated.CreateAmlAlertRequestAttributes{
			Amount:         details.Amount,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: rgenerated.CreateAmlAlertRequestRelationships{
			Balance: NewBalanceKey(details.BalanceID).AsRelation(),
		},
	}
}
func newAssetCreateRequest(id int64, details history2.CreateAssetRequest) *rgenerated.CreateAssetRequest {
	return &rgenerated.CreateAssetRequest{
		Key: rgenerated.NewKeyInt64(id, rgenerated.REQUEST_DETAILS_ASSET_CREATE),
		Attributes: rgenerated.CreateAssetRequestAttributes{
			Asset:                  details.Asset,
			Policies:               details.Policies,
			PreIssuanceAssetSigner: details.PreIssuedAssetSigner,
			MaxIssuanceAmount:      details.MaxIssuanceAmount,
			InitialPreissuedAmount: details.InitialPreissuedAmount,
			CreatorDetails:         details.CreatorDetails,
			Type:                   details.Type,
		},
	}
}
func newAssetUpdateRequest(id int64, details history2.UpdateAssetRequest) *rgenerated.UpdateAssetRequest {
	return &rgenerated.UpdateAssetRequest{
		Key: rgenerated.NewKeyInt64(id, rgenerated.REQUEST_DETAILS_ASSET_UPDATE),
		Attributes: rgenerated.UpdateAssetRequestAttributes{
			Policies:       details.Policies,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: rgenerated.UpdateAssetRequestRelationships{
			Asset: NewAssetKey(details.Asset).AsRelation(),
		},
	}
}
func newAtomicSwapRequest(id int64, details history2.CreateAtomicSwapRequest) *rgenerated.CreateAtomicSwapRequest {
	return &rgenerated.CreateAtomicSwapRequest{
		Key: rgenerated.NewKeyInt64(id, rgenerated.REQUEST_DETAILS_ATOMIC_SWAP),
		Attributes: rgenerated.CreateAtomicSwapRequestAttributes{
			BaseAmount: rgenerated.Amount(details.BaseAmount),
		},
		Relationships: rgenerated.CreateAtomicSwapRequestRelationships{
			Bid:        rgenerated.NewKeyInt64(int64(details.BidID), rgenerated.ASWAP_BID).AsRelation(),
			QuoteAsset: newQuoteAssetKey(details.QuoteAsset).AsRelation(),
		},
	}
}
func newAtomicSwapBidRequest(id int64, details history2.CreateAtomicSwapBidRequest) *rgenerated.CreateAtomicSwapBidRequest {
	quoteAssets := &rgenerated.RelationCollection{
		Data: make([]rgenerated.Key, 0, len(details.QuoteAssets)),
	}
	for _, quoteAsset := range details.QuoteAssets {
		quoteAssets.Data = append(quoteAssets.Data, newQuoteAssetKey(quoteAsset.Asset))
	}

	return &rgenerated.CreateAtomicSwapBidRequest{
		Key: rgenerated.NewKeyInt64(id, rgenerated.REQUEST_DETAILS_ASWAP_BID),
		Attributes: rgenerated.CreateAtomicSwapBidRequestAttributes{
			BaseAmount:     rgenerated.Amount(details.BaseAmount),
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: rgenerated.CreateAtomicSwapBidRequestRelationships{
			BaseBalance: NewBalanceKey(details.BaseBalance).AsRelation(),
			QuoteAssets: quoteAssets,
		},
	}
}
func newIssuanceRequest(id int64, details history2.CreateIssuanceRequest) *rgenerated.CreateIssuanceRequest {
	return &rgenerated.CreateIssuanceRequest{
		Key: rgenerated.NewKeyInt64(id, rgenerated.REQUEST_DETAILS_ISSUANCE),
		Attributes: rgenerated.CreateIssuanceRequestAttributes{
			Amount:         details.Amount,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: rgenerated.CreateIssuanceRequestRelationships{
			Asset:    NewAssetKey(details.Asset).AsRelation(),
			Receiver: NewBalanceKey(details.Receiver).AsRelation(),
		},
	}
}
func newLimitsUpdateRequest(id int64, details history2.UpdateLimitsRequest) *rgenerated.UpdateLimitsRequest {
	return &rgenerated.UpdateLimitsRequest{
		Key: rgenerated.NewKeyInt64(id, rgenerated.REQUEST_DETAILS_LIMITS_UPDATE),
		Attributes: rgenerated.UpdateLimitsRequestAttributes{
			CreatorDetails: details.CreatorDetails,
		},
	}
}
func newPreIssuanceRequest(id int64, details history2.CreatePreIssuanceRequest) *rgenerated.CreatePreIssuanceRequest {
	return &rgenerated.CreatePreIssuanceRequest{
		Key: rgenerated.NewKeyInt64(id, rgenerated.REQUEST_DETAILS_PRE_ISSUANCE),
		Attributes: rgenerated.CreatePreIssuanceRequestAttributes{
			Amount:    details.Amount,
			Signature: details.Signature,
			Reference: details.Reference,
		},
		Relationships: rgenerated.CreatePreIssuanceRequestRelationships{
			Asset: NewAssetKey(details.Asset).AsRelation(),
		},
	}
}
func newSaleRequest(id int64, details history2.CreateSaleRequest) *rgenerated.CreateSaleRequest {
	quoteAssets := &rgenerated.RelationCollection{
		Data: make([]rgenerated.Key, 0, len(details.QuoteAssets)),
	}
	for _, quoteAsset := range details.QuoteAssets {
		quoteAssets.Data = append(quoteAssets.Data, newQuoteAssetKey(quoteAsset.Asset))
	}

	return &rgenerated.CreateSaleRequest{
		Key: rgenerated.NewKeyInt64(id, rgenerated.REQUEST_DETAILS_SALE),
		Attributes: rgenerated.CreateSaleRequestAttributes{
			BaseAssetForHardCap: details.BaseAssetForHardCap,
			StartTime:           details.StartTime,
			EndTime:             details.EndTime,
			SaleType:            details.SaleType,
			CreatorDetails:      details.CreatorDetails,
		},
		Relationships: rgenerated.CreateSaleRequestRelationships{
			QuoteAssets:       quoteAssets,
			BaseAsset:         NewAssetKey(details.BaseAsset).AsRelation(),
			DefaultQuoteAsset: newQuoteAssetKey(details.DefaultQuoteAsset).AsRelation(),
		},
	}
}
func newChangeRoleRequest(id int64, details history2.ChangeRoleRequest) *rgenerated.ChangeRoleRequest {
	return &rgenerated.ChangeRoleRequest{
		Key: rgenerated.NewKeyInt64(id, rgenerated.REQUEST_DETAILS_CHANGE_ROLE),
		Attributes: rgenerated.ChangeRoleRequestAttributes{
			AccountRoleToSet: details.AccountRoleToSet,
			CreatorDetails:   details.CreatorDetails,
			SequenceNumber:   details.SequenceNumber,
		},
		Relationships: rgenerated.ChangeRoleRequestRelationships{
			AccountToUpdateRole: NewAccountKey(details.DestinationAccount).AsRelation(),
		},
	}
}
func newUpdateSaleDetailsRequest(id int64, details history2.UpdateSaleDetailsRequest) *rgenerated.UpdateSaleDetailsRequest {
	return &rgenerated.UpdateSaleDetailsRequest{
		Key: rgenerated.NewKeyInt64(id, rgenerated.REQUEST_DETAILS_UPDATE_SALE_DETAILS),
		Attributes: rgenerated.UpdateSaleDetailsRequestAttributes{
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: rgenerated.UpdateSaleDetailsRequestRelationships{
			Sale: NewSaleKey(int64(details.SaleID)).AsRelation(),
		},
	}
}
func newWithdrawalRequest(id int64, details history2.CreateWithdrawalRequest) *rgenerated.CreateWithdrawRequest {
	return &rgenerated.CreateWithdrawRequest{
		Key: rgenerated.NewKeyInt64(id, rgenerated.REQUEST_DETAILS_WITHDRAWAL),
		Attributes: rgenerated.CreateWithdrawRequestAttributes{
			Fee:            details.Fee,
			Amount:         details.Amount,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: rgenerated.CreateWithdrawRequestRelationships{
			Balance: NewBalanceKey(details.BalanceID).AsRelation(),
		},
	}
}

func newCreatePollRequest(id int64, details history2.CreatePollRequest) *regources.CreatePollRequest {
	return &regources.CreatePollRequest{
		Key: regources.NewKeyInt64(id, regources.TypeRequestDetailsCreatePoll),
		Attributes: regources.CreatePollRequestAttrs{
			PollData: regources.PollData{
				Type: details.PollData.Type,
			},
			StartTime:                details.StartTime,
			EndTime:                  details.EndTime,
			NumberOfChoices:          details.NumberOfChoices,
			PermissionType:           details.PermissionType,
			VoteConfirmationRequired: details.VoteConfirmationRequired,
			CreatorDetails:           details.CreatorDetails,
		},
		Relationships: regources.CreatePollRequestRelations{
			ResultProvider: NewAccountKey(details.ResultProviderID).AsRelation(),
		},
	}
}

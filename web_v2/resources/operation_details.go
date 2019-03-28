package resources

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/rgenerated"
)

//NewOperationDetails - populates operation details into appropriate resource
func NewOperationDetails(op history2.Operation) rgenerated.Resource {
	switch op.Type {
	case xdr.OperationTypeCreateAccount:
		return &rgenerated.CreateAccountOp{
			Key: rgenerated.NewKeyInt64(op.ID, rgenerated.OPERATIONS_CREATE_ACCOUNT),
			Relationships: rgenerated.CreateAccountOpRelationships{
				Account: NewAccountKey(op.Details.CreateAccount.AccountAddress).AsRelation(),
				Role:    NewAccountRoleKey(op.Details.CreateAccount.AccountRole).AsRelation(),
			},
		}
	case xdr.OperationTypeCreateIssuanceRequest:
		return newCreateIssuanceOpDetails(op.ID, *op.Details.CreateIssuanceRequest)
	case xdr.OperationTypeSetFees:
		return &rgenerated.SetFeeOp{
			Key:        rgenerated.NewKeyInt64(op.ID, rgenerated.OPERATIONS_SET_FEES),
			Attributes: rgenerated.SetFeeOpAttributes(*op.Details.SetFee),
		}
	case xdr.OperationTypeCreateWithdrawalRequest:
		return newCreateWithdrawalRequestOp(op.ID, *op.Details.CreateWithdrawRequest)
	case xdr.OperationTypeManageBalance:
		return newManageBalanceOp(op.ID, *op.Details.ManageBalance)
	case xdr.OperationTypeManageAsset:
		return newManageAssetOp(op.ID, *op.Details.ManageAsset)
	case xdr.OperationTypeCreatePreissuanceRequest:
		return newPreIssuanceRequestOp(op.ID, *op.Details.CreatePreIssuanceRequest)
	case xdr.OperationTypeManageAssetPair:
		return newManageAssetPairOp(op.ID, *op.Details.ManageAssetPair)
	case xdr.OperationTypeManageOffer:
		return newManageOfferOp(op.ID, *op.Details.ManageOffer)
	case xdr.OperationTypeManageInvoiceRequest:
		return rgenerated.NewKeyInt64(op.ID, rgenerated.OPERATIONS_MANAGE_INVOICE_REQUEST).GetKeyP()
	case xdr.OperationTypeReviewRequest:
		return newReviewRequestOp(op.ID, *op.Details.ReviewRequest)
	case xdr.OperationTypeCreateSaleRequest:
		return newCreateSaleRequestOp(op.ID, *op.Details.CreateSaleRequest)
	case xdr.OperationTypeCheckSaleState:
		return newCheckSaleStateOp(op.ID, *op.Details.CheckSaleState)
	case xdr.OperationTypeCreateAmlAlert:
		return newCreateAMLAlertRequestOp(op.ID, *op.Details.CreateAMLAlertRequest)
	case xdr.OperationTypeCreateChangeRoleRequest:
		return newChangeRoleRequestOp(op.ID, *op.Details.CreateChangeRoleRequest)
	case xdr.OperationTypePayment:
		return newPaymentOp(op.ID, *op.Details.Payment)
	case xdr.OperationTypeManageExternalSystemAccountIdPoolEntry:
		return newManageExternalSystemPool(op.ID, *op.Details.ManageExternalSystemPool)
	case xdr.OperationTypeBindExternalSystemAccountId:
		return &rgenerated.BindExternalSystemAccountIdOp{
			Key:        rgenerated.NewKeyInt64(op.ID, rgenerated.OPERATIONS_BIND_EXTERNAL_SYSTEM_ACCOUNT_ID),
			Attributes: rgenerated.BindExternalSystemAccountIdOpAttributes(*op.Details.BindExternalSystemAccount),
		}
	case xdr.OperationTypeManageSale:
		return &rgenerated.ManageSaleOp{
			Key:        rgenerated.NewKeyInt64(op.ID, rgenerated.OPERATIONS_MANAGE_SALE),
			Attributes: rgenerated.ManageSaleOpAttributes(*op.Details.ManageSale),
		}
	case xdr.OperationTypeManageKeyValue:
		return &rgenerated.ManageKeyValueOp{
			Key:        rgenerated.NewKeyInt64(op.ID, rgenerated.OPERATIONS_MANAGE_KEY_VALUE),
			Attributes: rgenerated.ManageKeyValueOpAttributes(*op.Details.ManageKeyValue),
		}
	case xdr.OperationTypeCreateManageLimitsRequest:
		return newCreateManageLimitsRequestOp(op.ID, *op.Details.CreateManageLimitsRequest)
	case xdr.OperationTypeManageContractRequest:
		return rgenerated.NewKeyInt64(op.ID, rgenerated.OPERATIONS_MANAGE_CONTRACT_REQUEST).GetKeyP()
	case xdr.OperationTypeManageContract:
		return rgenerated.NewKeyInt64(op.ID, rgenerated.OPERATIONS_MANAGE_CONTRACT).GetKeyP()
	case xdr.OperationTypeCancelSaleRequest:
		return rgenerated.NewKeyInt64(op.ID, rgenerated.OPERATIONS_CANCEL_SALE_REQUEST).GetKeyP()
	case xdr.OperationTypePayout:
		return newPayoutOp(op.ID, *op.Details.Payout)
	case xdr.OperationTypeManageAccountRole:
		return newManageAccountRole(op.ID, *op.Details.ManageAccountRole)
	case xdr.OperationTypeManageAccountRule:
		return newManageAccountRule(op.ID, *op.Details.ManageAccountRule)
	case xdr.OperationTypeCreateAswapBidRequest:
		return rgenerated.NewKeyInt64(op.ID, rgenerated.OPERATIONS_CREATE_ASWAP_BID_REQUEST).GetKeyP()
	case xdr.OperationTypeCancelAswapBid:
		return rgenerated.NewKeyInt64(op.ID, rgenerated.OPERATIONS_CANCEL_ASWAP_BID).GetKeyP()
	case xdr.OperationTypeCreateAswapRequest:
		return rgenerated.NewKeyInt64(op.ID, rgenerated.OPERATIONS_CREATE_ASWAP_REQUEST).GetKeyP()
	case xdr.OperationTypeManageSignerRole:
		return newManageSignerRole(op.ID, *op.Details.ManageSignerRole)
	case xdr.OperationTypeManageSignerRule:
		return newManageSignerRule(op.ID, *op.Details.ManageSignerRule)
	case xdr.OperationTypeManageSigner:
		return newManageSigner(op.ID, *op.Details.ManageSigner)
	case xdr.OperationTypeManageLimits:
		return newManageLimitsOp(op.ID, *op.Details.ManageLimits)
	case xdr.OperationTypeStamp:
		return newStampOp(op.ID, *op.Details.Stamp)
	case xdr.OperationTypeLicense:
		return newLicenseOp(op.ID, *op.Details.License)
	default:
		panic(errors.From(errors.New("unexpected operation type"), logan.F{
			"type": op.Type,
		}))
	}
}

// newManageLimitsOp - creates new instance of ManageLimitsOp
func newManageLimitsOp(id int64, details history2.ManageLimitsDetails) *rgenerated.ManageLimitsOp {
	result := rgenerated.ManageLimitsOp{
		Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_MANAGE_LIMITS),
		Attributes: rgenerated.ManageLimitsOpAttributes{
			Action: details.Action,
		},
	}

	switch details.Action {
	case xdr.ManageLimitsActionCreate:
		result.Attributes.Create = newManageLimitsCreationOp(*details.Creation)
	case xdr.ManageLimitsActionRemove:
		result.Attributes.Remove = &rgenerated.ManageLimitsRemovalOp{
			LimitsId: details.Removal.LimitsID,
		}
	default:
		panic(errors.From(errors.New("unexpected manage limits action"), logan.F{
			"action": details.Action,
		}))
	}

	return &result
}

// newManageLimitsCreationOp - creates new instance of ManageLimitsCreationOp
func newManageLimitsCreationOp(details history2.ManageLimitsCreationDetails) *rgenerated.ManageLimitsCreationOp {
	return &rgenerated.ManageLimitsCreationOp{
		AccountRole:     details.AccountRole,
		AccountAddress:  details.AccountAddress,
		StatsOpType:     details.StatsOpType,
		AssetCode:       details.AssetCode,
		IsConvertNeeded: details.IsConvertNeeded,
		DailyOut:        details.DailyOut,
		WeeklyOut:       details.WeeklyOut,
		AnnualOut:       details.AnnualOut,
		MonthlyOut:      details.MonthlyOut,
	}
}

// newReviewRequestOp - creates new instance of ReviewRequestOp
func newReviewRequestOp(id int64, details history2.ReviewRequestDetails) *rgenerated.ReviewRequestOp {
	return &rgenerated.ReviewRequestOp{
		Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_REVIEW_REQUEST),
		Attributes: rgenerated.ReviewRequestOpAttributes{
			Action:          details.Action,
			Reason:          details.Reason,
			RequestHash:     details.RequestHash,
			RequestId:       details.RequestID,
			IsFulfilled:     details.IsFulfilled,
			AddedTasks:      details.AddedTasks,
			RemovedTasks:    details.RemovedTasks,
			ExternalDetails: details.ExternalDetails,
		},
	}
}

// newManageExternalSystemPool - creates new instance of ManageExternalSystemPoolOp
func newManageExternalSystemPool(id int64, details history2.ManageExternalSystemPoolDetails) *rgenerated.ManageExternalSystemAccountIdPoolEntryOp {
	result := &rgenerated.ManageExternalSystemAccountIdPoolEntryOp{
		Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_MANAGE_EXTERNAL_SYSTEM_ACCOUNT_ID_POOL_ENTRY),
	}

	switch details.Action {
	case xdr.ManageExternalSystemAccountIdPoolEntryActionCreate:
		result.Attributes.Create = new(rgenerated.CreateExternalSystemPoolOp)
		*result.Attributes.Create = rgenerated.CreateExternalSystemPoolOp(*details.Create)
	case xdr.ManageExternalSystemAccountIdPoolEntryActionRemove:
		result.Attributes.Remove = new(rgenerated.RemoveExternalSystemPoolOp)
		*result.Attributes.Remove = rgenerated.RemoveExternalSystemPoolOp(*details.Remove)
	default:
		panic(errors.From(errors.New("unexpected action for manage ex sys id pool"), logan.F{
			"action": details.Action,
		}))
	}

	return result
}

// newChangeRoleRequest - creates new instance of CreateKYCRequestOp
func newChangeRoleRequestOp(id int64, details history2.CreateChangeRoleRequestDetails,
) *rgenerated.CreateChangeRoleRequestOp {
	return &rgenerated.CreateChangeRoleRequestOp{
		Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_CREATE_CHANGE_ROLE_REQUEST),
		Attributes: rgenerated.CreateChangeRoleRequestOpAttributes{
			CreatorDetails: details.CreatorDetails,
			AllTasks:       details.AllTasks,
		},
		Relationships: rgenerated.CreateChangeRoleRequestOpRelationships{
			AccountToUpdateRole: NewAccountKey(details.DestinationAccount).AsRelation(),
			Request:             NewRequestKey(details.RequestDetails.RequestID).AsRelation(),
			RoleToSet:           NewAccountRoleKey(details.AccountRoleToSet).AsRelation(),
		},
	}
}

// newCreateIssuanceOpDetails - creates new instance of CreateIssuanceRequestOp
func newCreateIssuanceOpDetails(id int64, details history2.CreateIssuanceRequestDetails) *rgenerated.CreateIssuanceRequestOp {
	return &rgenerated.CreateIssuanceRequestOp{
		Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_CREATE_ISSUANCE_REQUEST),
		Attributes: rgenerated.CreateIssuanceRequestOpAttributes{
			Fee:            details.Fee,
			Reference:      details.Reference,
			Amount:         details.Amount,
			CreatorDetails: details.CreatorDetails,
			AllTasks:       details.AllTasks,
		},
		Relationships: rgenerated.CreateIssuanceRequestOpRelationships{
			Asset:           NewAssetKey(details.Asset).AsRelation(),
			ReceiverAccount: NewAccountKey(details.ReceiverAccountAddress).AsRelation(),
			ReceiverBalance: NewBalanceKey(details.ReceiverBalanceAddress).AsRelation(),
			Request:         NewRequestKey(details.RequestDetails.RequestID).AsRelation(),
		},
	}
}

// newCreateWithdrawalRequestOp int6 creates new instance of
func newCreateWithdrawalRequestOp(id int64,
	details history2.CreateWithdrawRequestDetails) *rgenerated.CreateWithdrawRequestOp {
	return &rgenerated.CreateWithdrawRequestOp{
		Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_CREATE_WITHDRAWAL_REQUEST),
		Attributes: rgenerated.CreateWithdrawRequestOpAttributes{
			Amount:         details.Amount,
			Fee:            details.Fee,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: rgenerated.CreateWithdrawRequestOpRelationships{
			Balance: NewBalanceKey(details.BalanceAddress).AsRelation(),
		},
	}
}

// newManageBalanceOp - creates new instance of ManageBalanceOp
func newManageBalanceOp(id int64, details history2.ManageBalanceDetails) *rgenerated.ManageBalanceOp {
	return &rgenerated.ManageBalanceOp{
		Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_MANAGE_BALANCE),
		Attributes: rgenerated.ManageBalanceOpAttributes{
			Action:         details.Action,
			BalanceAddress: details.BalanceAddress,
		},
		Relationships: rgenerated.ManageBalanceOpRelationships{
			DestinationAccount: NewAccountKey(details.DestinationAccount).AsRelation(),
			Asset:              NewAssetKey(details.Asset).AsRelation(),
		},
	}
}

// newManageAssetOp - creates new instance of ManageAsset
func newManageAssetOp(id int64, details history2.ManageAssetDetails) *rgenerated.ManageAssetOp {
	return &rgenerated.ManageAssetOp{
		Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_MANAGE_ASSET),
		Attributes: rgenerated.ManageAssetOpAttributes{
			AssetCode:         details.AssetCode,
			Action:            details.Action,
			Policies:          details.Policies,
			CreatorDetails:    details.CreatorDetails,
			PreissuedSigner:   details.PreissuedSigner,
			MaxIssuanceAmount: details.MaxIssuanceAmount,
		},
		Relationships: rgenerated.ManageAssetOpRelationships{
			Request: NewRequestKey(details.RequestID).AsRelation(),
		},
	}
}

// newPreIssuanceRequestOp - creates new instance of CreatePreIssuanceRequestOp
func newPreIssuanceRequestOp(id int64, details history2.CreatePreIssuanceRequestDetails) *rgenerated.CreatePreIssuanceRequestOp {
	return &rgenerated.CreatePreIssuanceRequestOp{
		Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_CREATE_PREISSUANCE_REQUEST),
		Attributes: rgenerated.CreatePreIssuanceRequestOpAttributes{
			Amount: details.Amount,
		},
		Relationships: rgenerated.CreatePreIssuanceRequestOpRelationships{
			Asset:   NewAssetKey(details.AssetCode).AsRelation(),
			Request: NewRequestKey(details.RequestID).AsRelation(),
		},
	}
}

// newManageAssetPairOp - creates new instance of ManageAssetPairOp
func newManageAssetPairOp(id int64, details history2.ManageAssetPairDetails) *rgenerated.ManageAssetPairOp {
	return &rgenerated.ManageAssetPairOp{
		Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_MANAGE_ASSET_PAIR),
		Attributes: rgenerated.ManageAssetPairOpAttributes{
			PhysicalPrice:           details.PhysicalPrice,
			PhysicalPriceCorrection: details.PhysicalPriceCorrection,
			MaxPriceStep:            details.MaxPriceStep,
			Policies:                details.Policies,
		},
		Relationships: rgenerated.ManageAssetPairOpRelationships{
			BaseAsset:  NewAssetKey(details.BaseAsset).AsRelation(),
			QuoteAsset: NewAssetKey(details.QuoteAsset).AsRelation(),
		},
	}
}

// newManageOfferOp - creates new instance of ManageOfferOp
func newManageOfferOp(id int64, details history2.ManageOfferDetails) *rgenerated.ManageOfferOp {
	return &rgenerated.ManageOfferOp{
		Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_MANAGE_OFFER),
		Attributes: rgenerated.ManageOfferOpAttributes{
			OfferId:     details.OfferID,
			OrderBookId: details.OrderBookID,
			BaseAmount:  details.Amount,
			Price:       details.Price,
			IsBuy:       details.IsBuy,
			Fee:         details.Fee,
			IsDeleted:   details.IsDeleted,
		},
		Relationships: rgenerated.ManageOfferOpRelationships{
			BaseAsset:  NewAssetKey(details.BaseAsset).AsRelation(),
			QuoteAsset: NewAssetKey(details.QuoteAsset).AsRelation(),
		},
	}
}

// newCreateSaleRequestOp - creates new instance of CreateSaleRequestOp
func newCreateSaleRequestOp(id int64, details history2.CreateSaleRequestDetails) *rgenerated.CreateSaleRequestOp {
	quoteAssets := &rgenerated.RelationCollection{
		Data: make([]rgenerated.Key, 0, len(details.QuoteAssets)),
	}

	for _, quoteAsset := range details.QuoteAssets {
		quoteAssets.Data = append(quoteAssets.Data, newQuoteAssetKey(quoteAsset.Asset))
	}

	return &rgenerated.CreateSaleRequestOp{
		Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_CREATE_SALE_REQUEST),
		Attributes: rgenerated.CreateSaleRequestOpAttributes{
			StartTime:      details.StartTime,
			EndTime:        details.EndTime,
			SoftCap:        details.SoftCap,
			HardCap:        details.HardCap,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: rgenerated.CreateSaleRequestOpRelationships{
			QuoteAssets:       quoteAssets,
			Request:           NewRequestKey(details.RequestID).AsRelation(),
			BaseAsset:         NewAssetKey(details.BaseAsset).AsRelation(),
			DefaultQuoteAsset: NewAssetKey(details.DefaultQuoteAsset).AsRelation(),
		},
	}
}

// newCreateAMLAlertRequestOp - creates new instance of CreateAMLAlertRequestOp
func newCreateAMLAlertRequestOp(id int64, details history2.CreateAMLAlertRequestDetails) *rgenerated.CreateAmlAlertRequestOp {
	return &rgenerated.CreateAmlAlertRequestOp{
		Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_CREATE_AML_ALERT),
		Attributes: rgenerated.CreateAmlAlertRequestOpAttributes{
			Amount:         details.Amount,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: rgenerated.CreateAmlAlertRequestOpRelationships{
			Balance: NewBalanceKey(details.BalanceAddress).AsRelation(),
		},
	}
}

// newCheckSaleStateOp - creates new instance of CheckSaleStateOp
func newCheckSaleStateOp(id int64, details history2.CheckSaleStateDetails) *rgenerated.CheckSaleStateOp {
	return &rgenerated.CheckSaleStateOp{
		Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_CHECK_SALE_STATE),
		Attributes: rgenerated.CheckSaleStateOpAttributes{
			Effect: details.Effect,
		},
		Relationships: rgenerated.CheckSaleStateOpRelationships{
			Sale: NewSaleKey(details.SaleID).AsRelation(),
		},
	}

}

// newPaymentOp - creates new instance of PaymentOp
func newPaymentOp(id int64, details history2.PaymentDetails) *rgenerated.PaymentOp {
	return &rgenerated.PaymentOp{
		Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_PAYMENT_V2),
		Attributes: rgenerated.PaymentOpAttributes{
			Amount:                  details.Amount,
			SourceFee:               details.SourceFee,
			DestinationFee:          details.DestinationFee,
			SourcePayForDestination: details.SourcePayForDestination,
			Subject:                 details.Subject,
			Reference:               details.Reference,
		},
		Relationships: rgenerated.PaymentOpRelationships{
			AccountFrom: NewAccountKey(details.AccountFrom).AsRelation(),
			AccountTo:   NewAccountKey(details.AccountTo).AsRelation(),
			BalanceFrom: NewBalanceKey(details.BalanceFrom).AsRelation(),
			BalanceTo:   NewBalanceKey(details.BalanceTo).AsRelation(),
			Asset:       NewAssetKey(details.Asset).AsRelation(),
		},
	}

}

// newCreateManageLimitsRequestOp - creates new instance of CreateManageLimitsRequestOp
func newCreateManageLimitsRequestOp(id int64, details history2.CreateManageLimitsRequestDetails) *rgenerated.CreateManageLimitsRequestOp {
	return &rgenerated.CreateManageLimitsRequestOp{
		Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_CREATE_MANAGE_LIMITS_REQUEST),
		Attributes: rgenerated.CreateManageLimitsRequestOpAttributes{
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: rgenerated.CreateManageLimitsRequestOpRelationships{
			Request: NewRequestKey(details.RequestID).AsRelation(),
		},
	}
}

// newPayoutOp - creates new instance of PayoutOp
func newPayoutOp(id int64, details history2.PayoutDetails) *rgenerated.PayoutOp {
	return &rgenerated.PayoutOp{
		Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_PAYOUT),
		Attributes: rgenerated.PayoutOpAttributes{
			MaxPayoutAmount:      details.MaxPayoutAmount,
			MinAssetHolderAmount: details.MinAssetHolderAmount,
			MinPayoutAmount:      details.MinPayoutAmount,
			ExpectedFee:          details.ExpectedFee,
			ActualFee:            details.ActualFee,
			ActualPayoutAmount:   details.ActualPayoutAmount,
		},
		Relationships: rgenerated.PayoutOpRelationships{
			SourceAccount: NewAccountKey(details.SourceAccountAddress).AsRelation(),
			SourceBalance: NewBalanceKey(details.SourceBalanceAddress).AsRelation(),
			Asset:         NewAssetKey(details.Asset).AsRelation(),
		},
	}
}

func newLicenseOp(id int64, details history2.LicenseDetails) *rgenerated.LicenseOp {
	return &rgenerated.LicenseOp{
		Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_LICENSE),
		Attributes: rgenerated.LicenseOpAttributes{
			PrevLicenseHash: details.PrevLicenseHash,
			LedgerHash:      details.LedgerHash,
			DueDate:         details.DueDate,
			AdminCount:      details.AdminCount,
			Signatures:      details.Signatures,
		},
	}
}

func newStampOp(id int64, details history2.StampDetails) *rgenerated.StampOp {
	return &rgenerated.StampOp{
		Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_STAMP),
		Attributes: rgenerated.StampOpAttributes{
			LedgerHash:  details.LedgerHash,
			LicenseHash: details.LicenseHash,
		},
	}
}

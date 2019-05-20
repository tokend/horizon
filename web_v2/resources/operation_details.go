package resources

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
)

//NewOperationDetails - populates operation details into appropriate resource
func NewOperationDetails(op history2.Operation) regources.Resource {
	switch op.Type {
	case xdr.OperationTypeCreateAccount:
		return &regources.CreateAccountOp{
			Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_CREATE_ACCOUNT),
			Relationships: regources.CreateAccountOpRelationships{
				Account: NewAccountKey(op.Details.CreateAccount.AccountAddress).AsRelation(),
				Role:    NewAccountRoleKey(op.Details.CreateAccount.AccountRole).AsRelation(),
			},
		}
	case xdr.OperationTypeCreateIssuanceRequest:
		return newCreateIssuanceOpDetails(op.ID, *op.Details.CreateIssuanceRequest)
	case xdr.OperationTypeSetFees:
		return &regources.SetFeeOp{
			Key:        regources.NewKeyInt64(op.ID, regources.OPERATIONS_SET_FEES),
			Attributes: regources.SetFeeOpAttributes(*op.Details.SetFee),
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
		return regources.NewKeyInt64(op.ID, regources.OPERATIONS_MANAGE_INVOICE).GetKeyP()
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
		return &regources.BindExternalSystemAccountIdOp{
			Key:        regources.NewKeyInt64(op.ID, regources.OPERATIONS_BIND_EXTERNAL_SYSTEM_ACCOUNT_ID),
			Attributes: regources.BindExternalSystemAccountIdOpAttributes(*op.Details.BindExternalSystemAccount),
		}
	case xdr.OperationTypeManageSale:
		return &regources.ManageSaleOp{
			Key:        regources.NewKeyInt64(op.ID, regources.OPERATIONS_MANAGE_SALE),
			Attributes: regources.ManageSaleOpAttributes(*op.Details.ManageSale),
		}
	case xdr.OperationTypeManageKeyValue:
		return &regources.ManageKeyValueOp{
			Key:        regources.NewKeyInt64(op.ID, regources.OPERATIONS_MANAGE_KEY_VALUE),
			Attributes: regources.ManageKeyValueOpAttributes(*op.Details.ManageKeyValue),
		}
	case xdr.OperationTypeCreateManageLimitsRequest:
		return newCreateManageLimitsRequestOp(op.ID, *op.Details.CreateManageLimitsRequest)
	case xdr.OperationTypeManageContractRequest:
		return regources.NewKeyInt64(op.ID, regources.OPERATIONS_MANAGE_CONTRACT_REQUEST).GetKeyP()
	case xdr.OperationTypeManageContract:
		return regources.NewKeyInt64(op.ID, regources.OPERATIONS_MANAGE_CONTRACT).GetKeyP()
	case xdr.OperationTypeCancelSaleRequest:
		return regources.NewKeyInt64(op.ID, regources.OPERATIONS_CANCEL_SALE_REQUEST).GetKeyP()
	case xdr.OperationTypeCancelChangeRoleRequest:
		return regources.NewKeyInt64(op.ID, regources.OPERATIONS_CANCEL_CHANGE_ROLE_REQUEST).GetKeyP()
	case xdr.OperationTypePayout:
		return newPayoutOp(op.ID, *op.Details.Payout)
	case xdr.OperationTypeManageAccountRole:
		return newManageAccountRole(op.ID, *op.Details.ManageAccountRole)
	case xdr.OperationTypeManageAccountRule:
		return newManageAccountRule(op.ID, *op.Details.ManageAccountRule)
	case xdr.OperationTypeCreateAswapBidRequest:
		return regources.NewKeyInt64(op.ID, regources.OPERATIONS_CREATE_ASWAP_BID_REQUEST).GetKeyP()
	case xdr.OperationTypeCancelAswapBid:
		return regources.NewKeyInt64(op.ID, regources.OPERATIONS_CANCEL_ASWAP_BID).GetKeyP()
	case xdr.OperationTypeCreateAswapRequest:
		return regources.NewKeyInt64(op.ID, regources.OPERATIONS_CREATE_ASWAP_REQUEST).GetKeyP()
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
	case xdr.OperationTypeManageCreatePollRequest:
		return newManageCreatePollRequestOp(op.ID, *op.Details.ManageCreatePollRequest)
	case xdr.OperationTypeManagePoll:
		return newManagePollOp(op.ID, *op.Details.ManagePoll)
	case xdr.OperationTypeManageVote:
		return newManageVoteOp(op.ID, *op.Details.ManageVote)
	default:
		panic(errors.From(errors.New("unexpected operation type"), logan.F{
			"type": op.Type,
		}))
	}
}

func newManageCreatePollRequestOp(id int64, details history2.ManageCreatePollRequestDetails) *regources.ManageCreatePollRequestOp {
	manageCreateRequestPollOp := regources.ManageCreatePollRequestOp{
		Key: regources.NewKeyInt64(id, regources.OPERATIONS_MANAGE_CREATE_POLL_REQUEST),
		Attributes: &regources.ManageCreatePollRequestOpAttributes{
			Action: details.Action,
		},
	}

	switch details.Action {
	case xdr.ManageCreatePollRequestActionCreate:
		manageCreateRequestPollOp.Attributes.Create = &regources.CreatePollRequestOp{
			AllTasks:        details.CreateDetails.AllTasks,
			CreatorDetails:  details.CreateDetails.CreatorDetails,
			EndTime:         details.CreateDetails.EndTime,
			NumberOfChoices: uint64(details.CreateDetails.NumberOfChoices),
			PermissionType:  uint64(details.CreateDetails.PermissionType),
			PollData: regources.PollData{
				Type: details.CreateDetails.PollData.Type,
			},
			ResultProviderId:         details.CreateDetails.ResultProviderID,
			StartTime:                details.CreateDetails.StartTime,
			VoteConfirmationRequired: details.CreateDetails.VoteConfirmationRequired,
		}
		manageCreateRequestPollOp.Relationships = &regources.ManageCreatePollRequestOpRelationships{
			Request:        NewRequestKey(details.CreateDetails.RequestDetails.RequestID).AsRelation(),
			ResultProvider: NewAccountKey(details.CreateDetails.ResultProviderID).AsRelation(),
		}
	case xdr.ManageCreatePollRequestActionCancel:
		manageCreateRequestPollOp.Relationships = &regources.ManageCreatePollRequestOpRelationships{
			Request: NewRequestKey(details.CancelDetails.RequestID).AsRelation(),
		}
	default:
		panic(errors.From(errors.New("unexpected poll request action"), logan.F{
			"action": details.Action,
		}))
	}

	return &manageCreateRequestPollOp
}

func newManagePollOp(id int64, details history2.ManagePollDetails) *regources.ManagePollOp {
	managePollOp := regources.ManagePollOp{
		Key: regources.NewKeyInt64(id, regources.OPERATIONS_MANAGE_POLL),
		Attributes: regources.ManagePollOpAttributes{
			Action: details.Action,
			PollId: details.PollID,
		},
		Relationships: regources.ManagePollOpRelationships{
			Poll: NewPollKey(details.PollID).AsRelation(),
		},
	}

	switch details.Action {
	case xdr.ManagePollActionClose:
		managePollOp.Attributes.Close = &regources.ClosePollOp{
			Details:    regources.Details(details.ClosePoll.Details),
			PollId:     details.PollID,
			PollResult: details.ClosePoll.PollResult,
		}
	case xdr.ManagePollActionUpdateEndTime:
		managePollOp.Attributes.UpdateEndTime = &regources.UpdatePollEndTimeOp{
			NewEndTime: details.UpdatePollEndTime.EndTime,
		}
	case xdr.ManagePollActionCancel:
	default:
		panic(errors.From(errors.New("unexpected manage poll action"), logan.F{
			"action": details.Action,
		}))
	}

	return &managePollOp
}

func newManageVoteOp(id int64, details history2.ManageVoteDetails) *regources.ManageVoteOp {
	manageVoteOp := regources.ManageVoteOp{
		Key: regources.NewKeyInt64(id, regources.OPERATIONS_MANAGE_VOTE),
		Attributes: regources.ManageVoteOpAttributes{
			Action: details.Action,
		},
		Relationships: regources.ManageVoteOpRelationships{
			Poll: NewPollKey(details.PollID).AsRelation(),
		},
	}

	switch details.Action {
	case xdr.ManageVoteActionCreate:
		choice := uint64(details.VoteData.Single.Choice)
		manageVoteOp.Attributes.Create = &regources.CreateVoteOp{
			PollId: details.PollID,
		}
		if details.VoteData != nil {
			manageVoteOp.Attributes.Create.VoteData = regources.VoteData{
				PollType:     details.VoteData.PollType,
				SingleChoice: &choice,
			}
		}
	case xdr.ManageVoteActionRemove:
		manageVoteOp.Attributes.Remove = &regources.RemoveVoteOp{
			PollId: details.PollID,
		}
	default:
		panic(errors.From(errors.New("unexpected manage vote action"), logan.F{
			"action": details.Action,
		}))
	}

	return &manageVoteOp
}

// newManageLimitsOp - creates new instance of ManageLimitsOp
func newManageLimitsOp(id int64, details history2.ManageLimitsDetails) *regources.ManageLimitsOp {
	result := regources.ManageLimitsOp{
		Key: regources.NewKeyInt64(id, regources.OPERATIONS_MANAGE_LIMITS),
		Attributes: regources.ManageLimitsOpAttributes{
			Action: details.Action,
		},
	}

	switch details.Action {
	case xdr.ManageLimitsActionCreate:
		result.Attributes.Create = newManageLimitsCreationOp(*details.Creation)
	case xdr.ManageLimitsActionRemove:
		result.Attributes.Remove = &regources.ManageLimitsRemovalOp{
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
func newManageLimitsCreationOp(details history2.ManageLimitsCreationDetails) *regources.ManageLimitsCreationOp {
	return &regources.ManageLimitsCreationOp{
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
func newReviewRequestOp(id int64, details history2.ReviewRequestDetails) *regources.ReviewRequestOp {
	return &regources.ReviewRequestOp{
		Key: regources.NewKeyInt64(id, regources.OPERATIONS_REVIEW_REQUEST),
		Attributes: regources.ReviewRequestOpAttributes{
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
func newManageExternalSystemPool(id int64, details history2.ManageExternalSystemPoolDetails) *regources.ManageExternalSystemAccountIdPoolEntryOp {
	result := &regources.ManageExternalSystemAccountIdPoolEntryOp{
		Key: regources.NewKeyInt64(id, regources.OPERATIONS_MANAGE_EXTERNAL_SYSTEM_ACCOUNT_ID_POOL_ENTRY),
	}

	switch details.Action {
	case xdr.ManageExternalSystemAccountIdPoolEntryActionCreate:
		result.Attributes.Create = new(regources.CreateExternalSystemPoolOp)
		*result.Attributes.Create = regources.CreateExternalSystemPoolOp(*details.Create)
	case xdr.ManageExternalSystemAccountIdPoolEntryActionRemove:
		result.Attributes.Remove = new(regources.RemoveExternalSystemPoolOp)
		*result.Attributes.Remove = regources.RemoveExternalSystemPoolOp(*details.Remove)
	default:
		panic(errors.From(errors.New("unexpected action for manage ex sys id pool"), logan.F{
			"action": details.Action,
		}))
	}

	return result
}

// newChangeRoleRequest - creates new instance of CreateKYCRequestOp
func newChangeRoleRequestOp(id int64, details history2.CreateChangeRoleRequestDetails,
) *regources.CreateChangeRoleRequestOp {
	return &regources.CreateChangeRoleRequestOp{
		Key: regources.NewKeyInt64(id, regources.OPERATIONS_CREATE_CHANGE_ROLE_REQUEST),
		Attributes: regources.CreateChangeRoleRequestOpAttributes{
			CreatorDetails: details.CreatorDetails,
			AllTasks:       details.AllTasks,
		},
		Relationships: regources.CreateChangeRoleRequestOpRelationships{
			AccountToUpdateRole: NewAccountKey(details.DestinationAccount).AsRelation(),
			Request:             NewRequestKey(details.RequestDetails.RequestID).AsRelation(),
			RoleToSet:           NewAccountRoleKey(details.AccountRoleToSet).AsRelation(),
		},
	}
}

// newCreateIssuanceOpDetails - creates new instance of CreateIssuanceRequestOp
func newCreateIssuanceOpDetails(id int64, details history2.CreateIssuanceRequestDetails) *regources.CreateIssuanceRequestOp {
	return &regources.CreateIssuanceRequestOp{
		Key: regources.NewKeyInt64(id, regources.OPERATIONS_CREATE_ISSUANCE_REQUEST),
		Attributes: regources.CreateIssuanceRequestOpAttributes{
			Fee:            details.Fee,
			Reference:      details.Reference,
			Amount:         details.Amount,
			CreatorDetails: details.CreatorDetails,
			AllTasks:       details.AllTasks,
		},
		Relationships: regources.CreateIssuanceRequestOpRelationships{
			Asset:           NewAssetKey(details.Asset).AsRelation(),
			ReceiverAccount: NewAccountKey(details.ReceiverAccountAddress).AsRelation(),
			ReceiverBalance: NewBalanceKey(details.ReceiverBalanceAddress).AsRelation(),
			Request:         NewRequestKey(details.RequestDetails.RequestID).AsRelation(),
		},
	}
}

// newCreateWithdrawalRequestOp int6 creates new instance of
func newCreateWithdrawalRequestOp(id int64,
	details history2.CreateWithdrawRequestDetails) *regources.CreateWithdrawRequestOp {
	return &regources.CreateWithdrawRequestOp{
		Key: regources.NewKeyInt64(id, regources.OPERATIONS_CREATE_WITHDRAWAL_REQUEST),
		Attributes: regources.CreateWithdrawRequestOpAttributes{
			Amount:         details.Amount,
			Fee:            details.Fee,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.CreateWithdrawRequestOpRelationships{
			Balance: NewBalanceKey(details.BalanceAddress).AsRelation(),
		},
	}
}

// newManageBalanceOp - creates new instance of ManageBalanceOp
func newManageBalanceOp(id int64, details history2.ManageBalanceDetails) *regources.ManageBalanceOp {
	return &regources.ManageBalanceOp{
		Key: regources.NewKeyInt64(id, regources.OPERATIONS_MANAGE_BALANCE),
		Attributes: regources.ManageBalanceOpAttributes{
			Action:         details.Action,
			BalanceAddress: details.BalanceAddress,
		},
		Relationships: regources.ManageBalanceOpRelationships{
			DestinationAccount: NewAccountKey(details.DestinationAccount).AsRelation(),
			Asset:              NewAssetKey(details.Asset).AsRelation(),
		},
	}
}

// newManageAssetOp - creates new instance of ManageAsset
func newManageAssetOp(id int64, details history2.ManageAssetDetails) *regources.ManageAssetOp {
	return &regources.ManageAssetOp{
		Key: regources.NewKeyInt64(id, regources.OPERATIONS_MANAGE_ASSET),
		Attributes: regources.ManageAssetOpAttributes{
			AssetCode:         details.AssetCode,
			Action:            details.Action,
			Policies:          details.Policies,
			CreatorDetails:    details.CreatorDetails,
			PreIssuanceSigner: details.PreIssuanceSigner,
			MaxIssuanceAmount: details.MaxIssuanceAmount,
		},
		Relationships: regources.ManageAssetOpRelationships{
			Request: NewRequestKey(details.RequestID).AsRelation(),
		},
	}
}

// newPreIssuanceRequestOp - creates new instance of CreatePreIssuanceRequestOp
func newPreIssuanceRequestOp(id int64, details history2.CreatePreIssuanceRequestDetails) *regources.CreatePreIssuanceRequestOp {
	return &regources.CreatePreIssuanceRequestOp{
		Key: regources.NewKeyInt64(id, regources.OPERATIONS_CREATE_PREISSUANCE_REQUEST),
		Attributes: regources.CreatePreIssuanceRequestOpAttributes{
			Amount: details.Amount,
		},
		Relationships: regources.CreatePreIssuanceRequestOpRelationships{
			Asset:   NewAssetKey(details.AssetCode).AsRelation(),
			Request: NewRequestKey(details.RequestID).AsRelation(),
		},
	}
}

// newManageAssetPairOp - creates new instance of ManageAssetPairOp
func newManageAssetPairOp(id int64, details history2.ManageAssetPairDetails) *regources.ManageAssetPairOp {
	return &regources.ManageAssetPairOp{
		Key: regources.NewKeyInt64(id, regources.OPERATIONS_MANAGE_ASSET_PAIR),
		Attributes: regources.ManageAssetPairOpAttributes{
			PhysicalPrice:           details.PhysicalPrice,
			PhysicalPriceCorrection: details.PhysicalPriceCorrection,
			MaxPriceStep:            details.MaxPriceStep,
			Policies:                details.Policies,
		},
		Relationships: regources.ManageAssetPairOpRelationships{
			BaseAsset:  NewAssetKey(details.BaseAsset).AsRelation(),
			QuoteAsset: NewAssetKey(details.QuoteAsset).AsRelation(),
		},
	}
}

// newManageOfferOp - creates new instance of ManageOfferOp
func newManageOfferOp(id int64, details history2.ManageOfferDetails) *regources.ManageOfferOp {
	return &regources.ManageOfferOp{
		Key: regources.NewKeyInt64(id, regources.OPERATIONS_MANAGE_OFFER),
		Attributes: regources.ManageOfferOpAttributes{
			OfferId:     details.OfferID,
			OrderBookId: details.OrderBookID,
			BaseAmount:  details.Amount,
			Price:       details.Price,
			IsBuy:       details.IsBuy,
			Fee:         details.Fee,
			IsDeleted:   details.IsDeleted,
		},
		Relationships: regources.ManageOfferOpRelationships{
			BaseAsset:  NewAssetKey(details.BaseAsset).AsRelation(),
			QuoteAsset: NewAssetKey(details.QuoteAsset).AsRelation(),
		},
	}
}

// newCreateSaleRequestOp - creates new instance of CreateSaleRequestOp
func newCreateSaleRequestOp(id int64, details history2.CreateSaleRequestDetails) *regources.CreateSaleRequestOp {
	quoteAssets := &regources.RelationCollection{
		Data: make([]regources.Key, 0, len(details.QuoteAssets)),
	}

	for _, quoteAsset := range details.QuoteAssets {
		quoteAssets.Data = append(quoteAssets.Data, newQuoteAssetKey(quoteAsset.Asset))
	}

	return &regources.CreateSaleRequestOp{
		Key: regources.NewKeyInt64(id, regources.OPERATIONS_CREATE_SALE_REQUEST),
		Attributes: regources.CreateSaleRequestOpAttributes{
			StartTime:      details.StartTime,
			EndTime:        details.EndTime,
			SoftCap:        details.SoftCap,
			HardCap:        details.HardCap,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.CreateSaleRequestOpRelationships{
			QuoteAssets:       quoteAssets,
			Request:           NewRequestKey(details.RequestID).AsRelation(),
			BaseAsset:         NewAssetKey(details.BaseAsset).AsRelation(),
			DefaultQuoteAsset: NewAssetKey(details.DefaultQuoteAsset).AsRelation(),
		},
	}
}

// newCreateAMLAlertRequestOp - creates new instance of CreateAMLAlertRequestOp
func newCreateAMLAlertRequestOp(id int64, details history2.CreateAMLAlertRequestDetails) *regources.CreateAmlAlertRequestOp {
	return &regources.CreateAmlAlertRequestOp{
		Key: regources.NewKeyInt64(id, regources.OPERATIONS_CREATE_AML_ALERT),
		Attributes: regources.CreateAmlAlertRequestOpAttributes{
			Amount:         details.Amount,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.CreateAmlAlertRequestOpRelationships{
			Balance: NewBalanceKey(details.BalanceAddress).AsRelation(),
		},
	}
}

// newCheckSaleStateOp - creates new instance of CheckSaleStateOp
func newCheckSaleStateOp(id int64, details history2.CheckSaleStateDetails) *regources.CheckSaleStateOp {
	return &regources.CheckSaleStateOp{
		Key: regources.NewKeyInt64(id, regources.OPERATIONS_CHECK_SALE_STATE),
		Attributes: regources.CheckSaleStateOpAttributes{
			Effect: details.Effect,
		},
		Relationships: regources.CheckSaleStateOpRelationships{
			Sale: NewSaleKey(details.SaleID).AsRelation(),
		},
	}

}

// newPaymentOp - creates new instance of PaymentOp
func newPaymentOp(id int64, details history2.PaymentDetails) *regources.PaymentOp {
	return &regources.PaymentOp{
		Key: regources.NewKeyInt64(id, regources.OPERATIONS_PAYMENT),
		Attributes: regources.PaymentOpAttributes{
			Amount:                  details.Amount,
			SourceFee:               details.SourceFee,
			DestinationFee:          details.DestinationFee,
			SourcePayForDestination: details.SourcePayForDestination,
			Subject:                 details.Subject,
			Reference:               details.Reference,
		},
		Relationships: regources.PaymentOpRelationships{
			AccountFrom: NewAccountKey(details.AccountFrom).AsRelation(),
			AccountTo:   NewAccountKey(details.AccountTo).AsRelation(),
			BalanceFrom: NewBalanceKey(details.BalanceFrom).AsRelation(),
			BalanceTo:   NewBalanceKey(details.BalanceTo).AsRelation(),
			Asset:       NewAssetKey(details.Asset).AsRelation(),
		},
	}

}

// newCreateManageLimitsRequestOp - creates new instance of CreateManageLimitsRequestOp
func newCreateManageLimitsRequestOp(id int64, details history2.CreateManageLimitsRequestDetails) *regources.CreateManageLimitsRequestOp {
	return &regources.CreateManageLimitsRequestOp{
		Key: regources.NewKeyInt64(id, regources.OPERATIONS_CREATE_MANAGE_LIMITS_REQUEST),
		Attributes: regources.CreateManageLimitsRequestOpAttributes{
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.CreateManageLimitsRequestOpRelationships{
			Request: NewRequestKey(details.RequestID).AsRelation(),
		},
	}
}

// newPayoutOp - creates new instance of PayoutOp
func newPayoutOp(id int64, details history2.PayoutDetails) *regources.PayoutOp {
	return &regources.PayoutOp{
		Key: regources.NewKeyInt64(id, regources.OPERATIONS_PAYOUT),
		Attributes: regources.PayoutOpAttributes{
			MaxPayoutAmount:      details.MaxPayoutAmount,
			MinAssetHolderAmount: details.MinAssetHolderAmount,
			MinPayoutAmount:      details.MinPayoutAmount,
			ExpectedFee:          details.ExpectedFee,
			ActualFee:            details.ActualFee,
			ActualPayoutAmount:   details.ActualPayoutAmount,
		},
		Relationships: regources.PayoutOpRelationships{
			SourceAccount: NewAccountKey(details.SourceAccountAddress).AsRelation(),
			SourceBalance: NewBalanceKey(details.SourceBalanceAddress).AsRelation(),
			Asset:         NewAssetKey(details.Asset).AsRelation(),
		},
	}
}

func newLicenseOp(id int64, details history2.LicenseDetails) *regources.LicenseOp {
	return &regources.LicenseOp{
		Key: regources.NewKeyInt64(id, regources.OPERATIONS_LICENSE),
		Attributes: regources.LicenseOpAttributes{
			PrevLicenseHash: details.PrevLicenseHash,
			LedgerHash:      details.LedgerHash,
			DueDate:         details.DueDate,
			AdminCount:      details.AdminCount,
			Signatures:      details.Signatures,
		},
	}
}

func newStampOp(id int64, details history2.StampDetails) *regources.StampOp {
	return &regources.StampOp{
		Key: regources.NewKeyInt64(id, regources.OPERATIONS_STAMP),
		Attributes: regources.StampOpAttributes{
			LedgerHash:  details.LedgerHash,
			LicenseHash: details.LicenseHash,
		},
	}
}

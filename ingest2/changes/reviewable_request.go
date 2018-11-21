package changes

import (
	"encoding/hex"
	"encoding/json"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/history"
	"gitlab.com/tokend/horizon/utf8"
	"gitlab.com/tokend/regources"
)

type reviewableRequestChanges struct {
	storage reviewableRequestStorage
}

type reviewableRequestStorage interface {
	InsertReviewableRequest(request history.ReviewableRequest) error
	UpdateReviewableRequest(request history.ReviewableRequest) error
	ApproveReviewableRequest(id uint64) error
	RejectReviewableRequest(id uint64) error
	UpdateInvoices(contractID uint64, oldStates uint64, newStates uint64) error
}

func (c *reviewableRequestChanges) Created(lc LedgerChange) error {
	reviewableRequest := lc.LedgerChange.MustCreated().Data.MustReviewableRequest()
	histReviewableReq, err := c.convertReviewableRequest(&reviewableRequest, lc.LedgerCloseTime)
	if err != nil {
		return errors.Wrap(err, "failed to convert reviewable request", logan.F{
			"request":         reviewableRequest,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	err = c.storage.InsertReviewableRequest(*histReviewableReq)
	if err != nil {
		return errors.Wrap(err, "failed to insert reviewable request", logan.F{
			"request":         histReviewableReq,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	return nil
}

func (c *reviewableRequestChanges) Updated(lc LedgerChange) error {
	reviewableRequest := lc.LedgerChange.MustUpdated().Data.MustReviewableRequest()
	histReviewableRequest, err := c.convertReviewableRequest(&reviewableRequest, lc.LedgerCloseTime)
	if err != nil {
		return errors.Wrap(err, "failed to convert reviewable request", logan.F{
			"request":         reviewableRequest,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	err = c.storage.UpdateReviewableRequest(*histReviewableRequest)
	if err != nil {
		return errors.Wrap(err, "failed to update reviewable request", logan.F{
			"request":         histReviewableRequest,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	return nil
}

func (c *reviewableRequestChanges) Deleted(lc LedgerChange) error {
	key := lc.LedgerChange.MustRemoved().MustReviewableRequest()

	op := lc.Operation.Body.ReviewRequestOp
	switch op.Action {
	case xdr.ReviewRequestOpActionApprove:
		err := c.storage.ApproveReviewableRequest(uint64(key.RequestId))
		if err != nil {
			return errors.Wrap(err, "Failed to delete reviewable request", logan.F{
				"ledger_entry_key": key,
			})
		}
	case xdr.ReviewRequestOpActionPermanentReject:
		err := c.storage.RejectReviewableRequest(uint64(key.RequestId))
		if err != nil {
			return errors.Wrap(err, "Failed to delete reviewable request", logan.F{
				"ledger_entry_key": key,
			})
		}
	}

	return nil
}

func (c *reviewableRequestChanges) convertReviewableRequest(request *xdr.ReviewableRequestEntry, ledgerCloseTime time.Time) (*history.ReviewableRequest, error) {
	var reference *string
	if request.Reference != nil {
		reference = new(string)
		*reference = utf8.Scrub(string(*request.Reference))
	}

	details, err := c.getReviewableRequestDetails(&request.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get reviewable request details")
	}

	state := history.ReviewableRequestStatePending
	if string(request.RejectReason) != "" {
		state = history.ReviewableRequestStateRejected
	}

	result := history.ReviewableRequest{
		TotalOrderID: db2.TotalOrderID{
			ID: int64(request.RequestId),
		},
		Requestor:    request.Requestor.Address(),
		Reviewer:     request.Reviewer.Address(),
		Reference:    reference,
		RejectReason: string(request.RejectReason),
		RequestType:  request.Body.Type,
		RequestState: state,
		Hash:         hex.EncodeToString(request.Hash[:]),
		Details:      details,
		CreatedAt:    time.Unix(int64(request.CreatedAt), 0).UTC(),
		UpdatedAt:    ledgerCloseTime,
	}

	tasksExt, ok := request.Ext.GetTasksExt()
	if !ok {
		return &result, nil
	}

	result.AllTasks = uint32(tasksExt.AllTasks)
	result.PendingTasks = uint32(tasksExt.PendingTasks)

	var externalDetails []map[string]interface{}
	for _, item := range tasksExt.ExternalDetails {
		var comment map[string]interface{}
		_ = json.Unmarshal([]byte(item), &comment)
		externalDetails = append(externalDetails, comment)
	}

	// we use key "data" for compatibility with db2.Details
	// the value for the key "data" is a slice of map[string]interface{}
	result.ExternalDetails = map[string]interface{}{
		"data": externalDetails,
	}

	return &result, nil
}

func (c *reviewableRequestChanges) getAssetCreation(request *xdr.AssetCreationRequest) *history.AssetCreationRequest {
	var details map[string]interface{}
	// error is ignored on purpose
	_ = json.Unmarshal([]byte(request.Details), &details)
	return &history.AssetCreationRequest{
		Asset:                  string(request.Code),
		Policies:               int32(request.Policies),
		PreIssuedAssetSigner:   request.PreissuedAssetSigner.Address(),
		MaxIssuanceAmount:      amount.StringU(uint64(request.MaxIssuanceAmount)),
		InitialPreissuedAmount: amount.StringU(uint64(request.InitialPreissuedAmount)),
		Details:                details,
	}
}

func (c *reviewableRequestChanges) getAssetUpdate(request *xdr.AssetUpdateRequest) *history.AssetUpdateRequest {
	var details map[string]interface{}
	// error is ignored on purpose
	_ = json.Unmarshal([]byte(request.Details), &details)
	return &history.AssetUpdateRequest{
		Asset:    string(request.Code),
		Policies: int32(request.Policies),
		Details:  details,
	}
}

func (c *reviewableRequestChanges) getPreIssuanceRequest(request *xdr.PreIssuanceRequest) (*history.PreIssuanceRequest, error) {
	signature, err := xdr.MarshalBase64(request.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal signature")
	}

	return &history.PreIssuanceRequest{
		Asset:     string(request.Asset),
		Amount:    amount.StringU(uint64(request.Amount)),
		Signature: signature,
		Reference: string(request.Reference),
	}, nil
}

func (c *reviewableRequestChanges) getIssuanceRequest(request *xdr.IssuanceRequest) *history.IssuanceRequest {
	var details map[string]interface{}
	// error is ignored on purpose, we should not block ingest in case of such error
	_ = json.Unmarshal([]byte(request.ExternalDetails), &details)
	return &history.IssuanceRequest{
		Asset:           string(request.Asset),
		Amount:          amount.StringU(uint64(request.Amount)),
		Receiver:        request.Receiver.AsString(),
		ExternalDetails: details,
	}
}

func (c *reviewableRequestChanges) getWithdrawalRequest(request *xdr.WithdrawalRequest) *history.WithdrawalRequest {
	var details map[string]interface{}
	// error is ignored on purpose, we should not block ingest in case of such error
	_ = json.Unmarshal([]byte(request.ExternalDetails), &details)

	var preConfirmationDetails map[string]interface{}
	_ = json.Unmarshal([]byte(request.PreConfirmationDetails), &preConfirmationDetails)
	return &history.WithdrawalRequest{
		BalanceID:              request.Balance.AsString(),
		Amount:                 amount.StringU(uint64(request.Amount)),
		FixedFee:               amount.StringU(uint64(request.Fee.Fixed)),
		PercentFee:             amount.StringU(uint64(request.Fee.Percent)),
		ExternalDetails:        details,
		DestAssetCode:          string(request.Details.AutoConversion.DestAsset),
		DestAssetAmount:        amount.StringU(uint64(request.Details.AutoConversion.ExpectedAmount)),
		PreConfirmationDetails: preConfirmationDetails,
	}
}

func (c *reviewableRequestChanges) getAmlAlertRequest(request *xdr.AmlAlertRequest) *history.AmlAlertRequest {
	return &history.AmlAlertRequest{
		BalanceID: request.BalanceId.AsString(),
		Amount:    amount.StringU(uint64(request.Amount)),
		Reason:    string(request.Reason),
	}
}

func (c *reviewableRequestChanges) getSaleRequest(request *xdr.SaleCreationRequest) *history.SaleRequest {
	var quoteAssets []regources.SaleQuoteAsset
	for i := range request.QuoteAssets {
		quoteAssets = append(quoteAssets, regources.SaleQuoteAsset{
			Price:      regources.Amount(int64(request.QuoteAssets[i].Price)),
			QuoteAsset: string(request.QuoteAssets[i].QuoteAsset),
		})
	}

	var details map[string]interface{}
	// error is ignored on purpose, we should not block ingest in case of such error
	_ = json.Unmarshal([]byte(request.Details), &details)

	saleType := xdr.SaleTypeBasicSale
	baseAssetForHardCap := uint64(0)
	state := xdr.SaleStateNone
	switch request.Ext.V {
	case xdr.LedgerVersionEmptyVersion:
	case xdr.LedgerVersionTypedSale:
		saleType = request.Ext.MustSaleTypeExt().TypedSale.SaleType
	case xdr.LedgerVersionAllowToSpecifyRequiredBaseAssetAmountForHardCap:
		extV2 := request.Ext.MustExtV2()
		baseAssetForHardCap = uint64(extV2.RequiredBaseAssetForHardCap)
		saleType = extV2.SaleTypeExt.TypedSale.SaleType
	case xdr.LedgerVersionStatableSales:
		extV3 := request.Ext.MustExtV3()
		saleType = extV3.SaleTypeExt.TypedSale.SaleType
		baseAssetForHardCap = uint64(extV3.RequiredBaseAssetForHardCap)
		state = extV3.State
	default:
		panic(errors.Wrap(errors.New("Unexpected ledger version in getSaleRequest"),
			"failed to ingest sale request", logan.F{
				"actual_ledger_version": request.Ext.V.ShortString(),
			}))
	}

	return &history.SaleRequest{
		BaseAsset:           string(request.BaseAsset),
		DefaultQuoteAsset:   string(request.DefaultQuoteAsset),
		StartTime:           time.Unix(int64(request.StartTime), 0).UTC(),
		EndTime:             time.Unix(int64(request.EndTime), 0).UTC(),
		SoftCap:             amount.StringU(uint64(request.SoftCap)),
		HardCap:             amount.StringU(uint64(request.HardCap)),
		Details:             details,
		QuoteAssets:         quoteAssets,
		SaleType:            saleType,
		BaseAssetForHardCap: amount.StringU(baseAssetForHardCap),
		State:               state,
	}
}

func (c *reviewableRequestChanges) getLimitsUpdateRequest(request *xdr.LimitsUpdateRequest) *history.LimitsUpdateRequest {
	details, ok := request.Ext.GetDetails()
	var detailsMap map[string]interface{}
	if ok {
		limitsDetails := string(details)
		// error is ignored on purpose, we should not block ingest in case of such error
		_ = json.Unmarshal([]byte(limitsDetails), &detailsMap)
	}
	return &history.LimitsUpdateRequest{
		Details:      detailsMap,
		DocumentHash: hex.EncodeToString(request.DeprecatedDocumentHash[:]),
	}
}

func (c *reviewableRequestChanges) getPromotionUpdateRequest(request *xdr.PromotionUpdateRequest) *history.PromotionUpdateRequest {
	newPromorionData := c.getSaleRequest(&request.NewPromotionData)

	return &history.PromotionUpdateRequest{
		SaleID:           uint64(request.PromotionId),
		NewPromotionData: *newPromorionData,
	}
}

func (c *reviewableRequestChanges) getUpdateKYCRequest(request *xdr.UpdateKycRequest) *history.UpdateKYCRequest {
	var kycData map[string]interface{}
	// error is ignored on purpose, we should not block ingest in case of such error
	_ = json.Unmarshal([]byte(request.KycData), &kycData)

	var externalDetails []map[string]interface{}
	for _, item := range request.ExternalDetails {
		var comment map[string]interface{}
		_ = json.Unmarshal([]byte(item), &comment)
		externalDetails = append(externalDetails, comment)
	}

	return &history.UpdateKYCRequest{
		AccountToUpdateKYC: request.AccountToUpdateKyc.Address(),
		AccountTypeToSet:   request.AccountTypeToSet,
		KYCLevel:           uint32(request.KycLevel),
		KYCData:            kycData,
		AllTasks:           uint32(request.AllTasks),
		PendingTasks:       uint32(request.PendingTasks),
		SequenceNumber:     uint32(request.SequenceNumber),
		ExternalDetails:    externalDetails,
	}
}

func (c *reviewableRequestChanges) getUpdateSaleDetailsRequest(request *xdr.UpdateSaleDetailsRequest) *history.UpdateSaleDetailsRequest {
	var newDetails map[string]interface{}
	// error is ignored on purpose, we should not block ingest in case of such error
	_ = json.Unmarshal([]byte(request.NewDetails), &newDetails)

	return &history.UpdateSaleDetailsRequest{
		SaleID:     uint64(request.SaleId),
		NewDetails: newDetails,
	}
}

func (c *reviewableRequestChanges) getInvoiceRequest(request *xdr.InvoiceRequest) *history.InvoiceRequest {
	var details map[string]interface{}
	// error is ignored on purpose, we should not block ingest in case of such error
	_ = json.Unmarshal([]byte(request.Details), &details)

	var contractID *int64
	if request.ContractId != nil {
		tmpContractID := int64(*request.ContractId)
		contractID = &tmpContractID
	}

	return &history.InvoiceRequest{
		Asset:           string(request.Asset),
		Amount:          uint64(request.Amount),
		ContractID:      contractID,
		Details:         details,
		PayerBalance:    request.SenderBalance.AsString(),
		ReceiverBalance: request.ReceiverBalance.AsString(),
	}
}

func (c *reviewableRequestChanges) getContractRequest(request *xdr.ContractRequest) *history.ContractRequest {
	var details map[string]interface{}
	// error is ignored on purpose, we should not block ingest in case of such error
	_ = json.Unmarshal([]byte(request.Details), &details)

	return &history.ContractRequest{
		Escrow:    request.Escrow.Address(),
		Details:   details,
		StartTime: time.Unix(int64(request.StartTime), 0).UTC(),
		EndTime:   time.Unix(int64(request.EndTime), 0).UTC(),
	}
}

func (c *reviewableRequestChanges) getUpdateSaleEndTimeRequest(request *xdr.UpdateSaleEndTimeRequest) *history.UpdateSaleEndTimeRequest {
	return &history.UpdateSaleEndTimeRequest{
		SaleID:     uint64(request.SaleId),
		NewEndTime: time.Unix(int64(request.NewEndTime), 0).UTC(),
	}
}

func (c *reviewableRequestChanges) getAtomicSwapBidCreationRequest(request *xdr.ASwapBidCreationRequest,
) *history.AtomicSwapBidCreation {
	var details map[string]interface{}
	_ = json.Unmarshal([]byte(request.Details), &details)

	var quoteAssets []regources.AssetPrice
	for _, quoteAsset := range request.QuoteAssets {
		quoteAssets = append(quoteAssets, regources.AssetPrice{
			Asset: string(quoteAsset.QuoteAsset),
			Price: regources.Amount(quoteAsset.Price),
		})
	}

	return &history.AtomicSwapBidCreation{
		BaseBalance: request.BaseBalance.AsString(),
		BaseAmount:  uint64(request.Amount),
		Details:     details,
		QuoteAssets: quoteAssets,
	}
}

func (c *reviewableRequestChanges) getAtomicSwapRequest(request *xdr.ASwapRequest,
) *history.AtomicSwap {
	return &history.AtomicSwap{
		BidID:      uint64(request.BidId),
		BaseAmount: uint64(request.BaseAmount),
		QuoteAsset: string(request.QuoteAsset),
	}
}

func (c *reviewableRequestChanges) getReviewableRequestDetails(body *xdr.ReviewableRequestEntryBody) (history.ReviewableRequestDetails, error) {
	var details history.ReviewableRequestDetails
	var err error
	switch body.Type {
	case xdr.ReviewableRequestTypeAssetCreate:
		details.AssetCreation = c.getAssetCreation(body.AssetCreationRequest)
	case xdr.ReviewableRequestTypeAssetUpdate:
		details.AssetUpdate = c.getAssetUpdate(body.AssetUpdateRequest)
	case xdr.ReviewableRequestTypeIssuanceCreate:
		details.IssuanceCreate = c.getIssuanceRequest(body.IssuanceRequest)
	case xdr.ReviewableRequestTypePreIssuanceCreate:
		details.PreIssuanceCreate, err = c.getPreIssuanceRequest(body.PreIssuanceRequest)
		if err != nil {
			return details, errors.Wrap(err, "failed to get pre issuance request")
		}
	case xdr.ReviewableRequestTypeWithdraw:
		details.Withdraw = c.getWithdrawalRequest(body.WithdrawalRequest)
	case xdr.ReviewableRequestTypeSale:
		details.Sale = c.getSaleRequest(body.SaleCreationRequest)
	case xdr.ReviewableRequestTypeLimitsUpdate:
		details.LimitsUpdate = c.getLimitsUpdateRequest(body.LimitsUpdateRequest)
	case xdr.ReviewableRequestTypeTwoStepWithdrawal:
		details.TwoStepWithdraw = c.getWithdrawalRequest(body.TwoStepWithdrawalRequest)
	case xdr.ReviewableRequestTypeAmlAlert:
		details.AmlAlert = c.getAmlAlertRequest(body.AmlAlertRequest)
	case xdr.ReviewableRequestTypeUpdateKyc:
		details.UpdateKYC = c.getUpdateKYCRequest(body.UpdateKycRequest)
	case xdr.ReviewableRequestTypeUpdateSaleDetails:
		details.UpdateSaleDetails = c.getUpdateSaleDetailsRequest(body.UpdateSaleDetailsRequest)
	case xdr.ReviewableRequestTypeInvoice:
		details.Invoice = c.getInvoiceRequest(body.InvoiceRequest)
	case xdr.ReviewableRequestTypeContract:
		details.Contract = c.getContractRequest(body.ContractRequest)
	case xdr.ReviewableRequestTypeUpdateSaleEndTime:
		details.UpdateSaleEndTimeRequest = c.getUpdateSaleEndTimeRequest(body.UpdateSaleEndTimeRequest)
	case xdr.ReviewableRequestTypeUpdatePromotion:
		details.PromotionUpdate = c.getPromotionUpdateRequest(body.PromotionUpdateRequest)
	case xdr.ReviewableRequestTypeCreateAtomicSwapBid:
		details.AtomicSwapBidCreation = c.getAtomicSwapBidCreationRequest(body.ASwapBidCreationRequest)
	case xdr.ReviewableRequestTypeAtomicSwap:
		details.AtomicSwap = c.getAtomicSwapRequest(body.ASwapRequest)
	default:
		return details, errors.From(errors.New("unexpected reviewable request type"), map[string]interface{}{
			"request_type": body.Type.String(),
		})
	}

	return details, nil
}

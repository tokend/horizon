package ingest

import (
	"encoding/hex"
	"encoding/json"
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/utf8"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

func reviewableRequestCreate(is *Session, ledgerEntry *xdr.LedgerEntry) error {
	reviewableRequest := ledgerEntry.Data.ReviewableRequest
	if reviewableRequest == nil {
		return errors.New("expected reviewable request not to be nil")
	}

	histReviewableRequest, err := convertReviewableRequest(reviewableRequest, is.Cursor.LedgerCloseTime())
	if err != nil {
		return errors.Wrap(err, "failed to convert reviewable request")
	}

	err = is.Ingestion.HistoryQ().ReviewableRequests().Insert(*histReviewableRequest)
	if err != nil {
		return errors.Wrap(err, "failed to create reviewable request")
	}

	return nil
}

func reviewableRequestUpdate(is *Session, ledgerEntry *xdr.LedgerEntry) error {
	reviewableRequest := ledgerEntry.Data.ReviewableRequest
	if reviewableRequest == nil {
		return errors.New("expected reviewable request not to be nil")
	}

	histReviewableRequest, err := convertReviewableRequest(reviewableRequest, is.Cursor.LedgerCloseTime())
	if err != nil {
		return errors.Wrap(err, "failed to convert reviewable request")
	}

	err = is.Ingestion.HistoryQ().ReviewableRequests().Update(*histReviewableRequest)
	if err != nil {
		return errors.Wrap(err, "failed to update reviewable request")
	}

	return nil
}

func reviewableRequestDelete(is *Session, key *xdr.LedgerKey) error {
	requestKey := key.ReviewableRequest
	if requestKey == nil {
		return errors.New("expected reviewable request key not to be nil")
	}

	// approve it since the request is most likely to be auto-reviewed
	// the case when it's a permanent reject will be handled later in ingest operation
	err := is.Ingestion.HistoryQ().ReviewableRequests().Approve(uint64(requestKey.RequestId))

	if err != nil {
		return errors.Wrap(err, "Failed to delete reviewable request")
	}

	return nil
}

func convertReviewableRequest(request *xdr.ReviewableRequestEntry, ledgerCloseTime time.Time) (*history.ReviewableRequest, error) {
	var reference *string
	if request.Reference != nil {
		reference = new(string)
		*reference = utf8.Scrub(string(*request.Reference))
	}

	details, err := getReviewableRequestDetails(&request.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get reviewable request details")
	}

	state := history.ReviewableRequestStatePending
	if string(request.RejectReason) != "" {
		state = history.ReviewableRequestStateRejected
	}

	return &history.ReviewableRequest{
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
	}, nil
}

func getAssetCreation(request *xdr.AssetCreationRequest) *history.AssetCreationRequest {
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

func getAssetUpdate(request *xdr.AssetUpdateRequest) *history.AssetUpdateRequest {
	var details map[string]interface{}
	// error is ignored on purpose
	_ = json.Unmarshal([]byte(request.Details), &details)
	return &history.AssetUpdateRequest{
		Asset:    string(request.Code),
		Policies: int32(request.Policies),
		Details:  details,
	}
}

func getPreIssuanceRequest(request *xdr.PreIssuanceRequest) (*history.PreIssuanceRequest, error) {
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

func getIssuanceRequest(request *xdr.IssuanceRequest) *history.IssuanceRequest {
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

func getWithdrawalRequest(request *xdr.WithdrawalRequest) *history.WithdrawalRequest {
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

func getAmlAlertRequest(request *xdr.AmlAlertRequest) *history.AmlAlertRequest {
	return &history.AmlAlertRequest{
		BalanceID: request.BalanceId.AsString(),
		Amount:    amount.StringU(uint64(request.Amount)),
		Reason:    string(request.Reason),
	}
}

func getSaleRequest(request *xdr.SaleCreationRequest) *history.SaleRequest {
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
	if request.Ext.SaleTypeExt != nil {
		saleType = request.Ext.SaleTypeExt.TypedSale.SaleType
	}

	var baseAssetForHardCap uint64 = 0
	if extV2, ok := request.Ext.GetExtV2(); ok {
		baseAssetForHardCap = uint64(extV2.RequiredBaseAssetForHardCap)
		saleType = extV2.SaleTypeExt.TypedSale.SaleType
	}

	state := xdr.SaleStateNone
	if request.Ext.ExtV3 != nil {
		state = request.Ext.ExtV3.State
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

func getLimitsUpdateRequest(request *xdr.LimitsUpdateRequest) *history.LimitsUpdateRequest {
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

func getPromotionUpdateRequest(request *xdr.PromotionUpdateRequest) *history.PromotionUpdateRequest {
	newPromorionData := getSaleRequest(&request.NewPromotionData)

	return &history.PromotionUpdateRequest{
		SaleID:           uint64(request.PromotionId),
		NewPromotionData: *newPromorionData,
	}
}

func getUpdateKYCRequest(request *xdr.UpdateKycRequest) *history.UpdateKYCRequest {
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

func getUpdateSaleDetailsRequest(request *xdr.UpdateSaleDetailsRequest) *history.UpdateSaleDetailsRequest {
	var newDetails map[string]interface{}
	// error is ignored on purpose, we should not block ingest in case of such error
	_ = json.Unmarshal([]byte(request.NewDetails), &newDetails)

	return &history.UpdateSaleDetailsRequest{
		SaleID:     uint64(request.SaleId),
		NewDetails: newDetails,
	}
}

func getInvoiceRequest(request *xdr.InvoiceRequest) *history.InvoiceRequest {
	var details map[string]interface{}
	// error is ignored on purpose, we should not block ingest in case of such error
	_ = json.Unmarshal([]byte(request.Details), &details)

	var contractID *int64
	if request.ContractId != nil {
		tmpContractID := int64(*request.ContractId)
		contractID = &tmpContractID
	}

	return &history.InvoiceRequest{
		Asset:      string(request.Asset),
		Amount:     uint64(request.Amount),
		ContractID: contractID,
		Details:    details,
	}
}

func getContractRequest(request *xdr.ContractRequest) *history.ContractRequest {
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

func getUpdateSaleEndTimeRequest(request *xdr.UpdateSaleEndTimeRequest) *history.UpdateSaleEndTimeRequest {
	return &history.UpdateSaleEndTimeRequest{
		SaleID:     uint64(request.SaleId),
		NewEndTime: time.Unix(int64(request.NewEndTime), 0).UTC(),
	}
}

func getReviewableRequestDetails(body *xdr.ReviewableRequestEntryBody) (history.ReviewableRequestDetails, error) {
	var details history.ReviewableRequestDetails
	var err error
	switch body.Type {
	case xdr.ReviewableRequestTypeAssetCreate:
		details.AssetCreation = getAssetCreation(body.AssetCreationRequest)
	case xdr.ReviewableRequestTypeAssetUpdate:
		details.AssetUpdate = getAssetUpdate(body.AssetUpdateRequest)
	case xdr.ReviewableRequestTypeIssuanceCreate:
		details.IssuanceCreate = getIssuanceRequest(body.IssuanceRequest)
	case xdr.ReviewableRequestTypePreIssuanceCreate:
		details.PreIssuanceCreate, err = getPreIssuanceRequest(body.PreIssuanceRequest)
		if err != nil {
			return details, errors.Wrap(err, "failed to get pre issuance request")
		}
	case xdr.ReviewableRequestTypeWithdraw:
		details.Withdraw = getWithdrawalRequest(body.WithdrawalRequest)
	case xdr.ReviewableRequestTypeSale:
		details.Sale = getSaleRequest(body.SaleCreationRequest)
	case xdr.ReviewableRequestTypeLimitsUpdate:
		details.LimitsUpdate = getLimitsUpdateRequest(body.LimitsUpdateRequest)
	case xdr.ReviewableRequestTypeTwoStepWithdrawal:
		details.TwoStepWithdraw = getWithdrawalRequest(body.TwoStepWithdrawalRequest)
	case xdr.ReviewableRequestTypeAmlAlert:
		details.AmlAlert = getAmlAlertRequest(body.AmlAlertRequest)
	case xdr.ReviewableRequestTypeUpdateKyc:
		details.UpdateKYC = getUpdateKYCRequest(body.UpdateKycRequest)
	case xdr.ReviewableRequestTypeUpdateSaleDetails:
		details.UpdateSaleDetails = getUpdateSaleDetailsRequest(body.UpdateSaleDetailsRequest)
	case xdr.ReviewableRequestTypeInvoice:
		details.Invoice = getInvoiceRequest(body.InvoiceRequest)
	case xdr.ReviewableRequestTypeContract:
		details.Contract = getContractRequest(body.ContractRequest)
	case xdr.ReviewableRequestTypeUpdateSaleEndTime:
		details.UpdateSaleEndTimeRequest = getUpdateSaleEndTimeRequest(body.UpdateSaleEndTimeRequest)
	case xdr.ReviewableRequestTypeUpdatePromotion:
		details.PromotionUpdate = getPromotionUpdateRequest(body.PromotionUpdateRequest)
	default:
		return details, errors.From(errors.New("unexpected reviewable request type"), map[string]interface{}{
			"request_type": body.Type.String(),
		})
	}

	return details, nil
}

package ingest

import (
	"encoding/hex"
	"encoding/json"
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/history"
	"gitlab.com/tokend/horizon/utf8"
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

	tasksExt := request.Tasks

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

func getAssetCreation(request *xdr.AssetCreationRequest) *history.AssetCreationRequest {
	var details map[string]interface{}
	// error is ignored on purpose
	_ = json.Unmarshal([]byte(request.CreatorDetails), &details)
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
	_ = json.Unmarshal([]byte(request.CreatorDetails), &details)
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
	_ = json.Unmarshal([]byte(request.CreatorDetails), &details)

	return &history.WithdrawalRequest{
		BalanceID:       request.Balance.AsString(),
		Amount:          amount.StringU(uint64(request.Amount)),
		FixedFee:        amount.StringU(uint64(request.Fee.Fixed)),
		PercentFee:      amount.StringU(uint64(request.Fee.Percent)),
		ExternalDetails: details,
	}
}

func getAmlAlertRequest(request *xdr.AmlAlertRequest) *history.AmlAlertRequest {
	return &history.AmlAlertRequest{
		BalanceID: request.BalanceId.AsString(),
		Amount:    amount.StringU(uint64(request.Amount)),
		Reason:    string(request.CreatorDetails),
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
	_ = json.Unmarshal([]byte(request.CreatorDetails), &details)

	saleType := request.SaleTypeExt.SaleType
	baseAssetForHardCap := uint64(request.RequiredBaseAssetForHardCap)

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
	}
}

func getLimitsUpdateRequest(request *xdr.LimitsUpdateRequest) *history.LimitsUpdateRequest {
	var detailsMap map[string]interface{}
	limitsDetails := string(request.CreatorDetails)
	// error is ignored on purpose, we should not block ingest in case of such error
	_ = json.Unmarshal([]byte(limitsDetails), &detailsMap)
	return &history.LimitsUpdateRequest{
		Details:      detailsMap,
		DocumentHash: hex.EncodeToString(request.DeprecatedDocumentHash[:]),
	}
}

func getUpdateKYCRequest(request *xdr.ChangeRoleRequest) *history.ChangeRoleRequest {
	var kycData map[string]interface{}
	// error is ignored on purpose, we should not block ingest in case of such error
	_ = json.Unmarshal([]byte(request.KycData), &kycData)

	return &history.ChangeRoleRequest{
		DestinationAccount: request.DestinationAccount.Address(),
		AccountRoleToSet:   uint64(request.AccountRoleToSet),
		KYCData:            kycData,
		SequenceNumber:     uint32(request.SequenceNumber),
	}
}

func getUpdateSaleDetailsRequest(request *xdr.UpdateSaleDetailsRequest) *history.UpdateSaleDetailsRequest {
	var newDetails map[string]interface{}
	// error is ignored on purpose, we should not block ingest in case of such error
	_ = json.Unmarshal([]byte(request.CreatorDetails), &newDetails)

	return &history.UpdateSaleDetailsRequest{
		SaleID:     uint64(request.SaleId),
		NewDetails: newDetails,
	}
}

func getInvoiceRequest(request *xdr.InvoiceRequest) *history.InvoiceRequest {
	var details map[string]interface{}
	// error is ignored on purpose, we should not block ingest in case of such error
	_ = json.Unmarshal([]byte(request.CreatorDetails), &details)

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

func getContractRequest(request *xdr.ContractRequest) *history.ContractRequest {
	var details map[string]interface{}
	// error is ignored on purpose, we should not block ingest in case of such error
	_ = json.Unmarshal([]byte(request.CreatorDetails), &details)

	return &history.ContractRequest{
		Escrow:    request.Escrow.Address(),
		Details:   details,
		StartTime: time.Unix(int64(request.StartTime), 0).UTC(),
		EndTime:   time.Unix(int64(request.EndTime), 0).UTC(),
	}
}

func getAtomicSwapBidCreationRequest(request *xdr.ASwapBidCreationRequest,
) *history.AtomicSwapBidCreation {
	var details map[string]interface{}
	_ = json.Unmarshal([]byte(request.CreatorDetails), &details)

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

func getAtomicSwapRequest(request *xdr.ASwapRequest,
) *history.AtomicSwap {
	return &history.AtomicSwap{
		BidID:      uint64(request.BidId),
		BaseAmount: uint64(request.BaseAmount),
		QuoteAsset: string(request.QuoteAsset),
	}
}

func getReviewableRequestDetails(body *xdr.ReviewableRequestEntryBody) (history.ReviewableRequestDetails, error) {
	var details history.ReviewableRequestDetails
	var err error
	switch body.Type {
	case xdr.ReviewableRequestTypeCreateAsset:
		details.AssetCreation = getAssetCreation(body.AssetCreationRequest)
	case xdr.ReviewableRequestTypeUpdateAsset:
		details.AssetUpdate = getAssetUpdate(body.AssetUpdateRequest)
	case xdr.ReviewableRequestTypeCreateIssuance:
		details.IssuanceCreate = getIssuanceRequest(body.IssuanceRequest)
	case xdr.ReviewableRequestTypeCreatePreIssuance:
		details.PreIssuanceCreate, err = getPreIssuanceRequest(body.PreIssuanceRequest)
		if err != nil {
			return details, errors.Wrap(err, "failed to get pre issuance request")
		}
	case xdr.ReviewableRequestTypeCreateWithdraw:
		details.Withdraw = getWithdrawalRequest(body.WithdrawalRequest)
	case xdr.ReviewableRequestTypeCreateSale:
		details.Sale = getSaleRequest(body.SaleCreationRequest)
	case xdr.ReviewableRequestTypeUpdateLimits:
		details.LimitsUpdate = getLimitsUpdateRequest(body.LimitsUpdateRequest)
	case xdr.ReviewableRequestTypeCreateAmlAlert:
		details.AmlAlert = getAmlAlertRequest(body.AmlAlertRequest)
	case xdr.ReviewableRequestTypeChangeRole:
		details.ChangeRole = getUpdateKYCRequest(body.ChangeRoleRequest)
	case xdr.ReviewableRequestTypeUpdateSaleDetails:
		details.UpdateSaleDetails = getUpdateSaleDetailsRequest(body.UpdateSaleDetailsRequest)
	case xdr.ReviewableRequestTypeCreateInvoice:
		details.Invoice = getInvoiceRequest(body.InvoiceRequest)
	case xdr.ReviewableRequestTypeManageContract:
		details.Contract = getContractRequest(body.ContractRequest)
	case xdr.ReviewableRequestTypeCreateAtomicSwapBid:
		details.AtomicSwapBidCreation = getAtomicSwapBidCreationRequest(body.ASwapBidCreationRequest)
	case xdr.ReviewableRequestTypeCreateAtomicSwap:
		details.AtomicSwap = getAtomicSwapRequest(body.ASwapRequest)
	default:
		return details, errors.From(errors.New("unexpected reviewable request type"), map[string]interface{}{
			"request_type": body.Type.String(),
		})
	}

	return details, nil
}

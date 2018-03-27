package ingest

import (
	"encoding/hex"
	"encoding/json"
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/utf8"
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

func getAssetCreation(request *xdr.AssetCreationRequest) history.AssetCreationRequest {
	var details map[string]interface{}
	// error is ignored on purpose
	_ = json.Unmarshal([]byte(request.Details), &details)
	return history.AssetCreationRequest{
		Asset:                  string(request.Code),
		Policies:               int32(request.Policies),
		PreIssuedAssetSigner:   request.PreissuedAssetSigner.Address(),
		MaxIssuanceAmount:      amount.StringU(uint64(request.MaxIssuanceAmount)),
		InitialPreissuedAmount: amount.StringU(uint64(request.InitialPreissuedAmount)),
		Details:                details,
	}
}

func getAssetUpdate(request *xdr.AssetUpdateRequest) history.AssetUpdateRequest {
	var details map[string]interface{}
	// error is ignored on purpose
	_ = json.Unmarshal([]byte(request.Details), &details)
	return history.AssetUpdateRequest{
		Asset:    string(request.Code),
		Policies: int32(request.Policies),
		Details:  details,
	}
}

func getPreIssuanceRequest(request *xdr.PreIssuanceRequest) (history.PreIssuanceRequest, error) {
	signature, err := xdr.MarshalBase64(request.Signature)
	if err != nil {
		return history.PreIssuanceRequest{}, errors.Wrap(err, "failed to marshal signature")
	}

	return history.PreIssuanceRequest{
		Asset:     string(request.Asset),
		Amount:    amount.StringU(uint64(request.Amount)),
		Signature: signature,
		Reference: string(request.Reference),
	}, nil
}

func getIssuanceRequest(request *xdr.IssuanceRequest) history.IssuanceRequest {
	var details map[string]interface{}
	// error is ignored on purpose, we should not block ingest in case of such error
	_ = json.Unmarshal([]byte(request.ExternalDetails), &details)
	return history.IssuanceRequest{
		Asset:           string(request.Asset),
		Amount:          amount.StringU(uint64(request.Amount)),
		Receiver:        request.Receiver.AsString(),
		ExternalDetails: details,
	}
}

func getWithdrawalRequest(request *xdr.WithdrawalRequest) history.WithdrawalRequest {
	var details map[string]interface{}
	// error is ignored on purpose, we should not block ingest in case of such error
	_ = json.Unmarshal([]byte(request.ExternalDetails), &details)

	var preConfirmationDetails map[string]interface{}
	_ = json.Unmarshal([]byte(request.PreConfirmationDetails), &preConfirmationDetails)
	return history.WithdrawalRequest{
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

func getAmlAlertRequest(request *xdr.AmlAlertRequest) history.AmlAlertRequest {
	return history.AmlAlertRequest{
		BalanceID: request.BalanceId.AsString(),
		Amount:    amount.StringU(uint64(request.Amount)),
		Reason:    string(request.Reason),
	}
}

func getSaleRequest(request *xdr.SaleCreationRequest) history.SaleRequest {
	var quoteAssets []history.SaleQuoteAsset
	for i := range request.QuoteAssets {
		quoteAssets = append(quoteAssets, history.SaleQuoteAsset{
			Price:      amount.StringU(uint64(request.QuoteAssets[i].Price)),
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

	return history.SaleRequest{
		BaseAsset:         string(request.BaseAsset),
		DefaultQuoteAsset: string(request.DefaultQuoteAsset),
		StartTime:         time.Unix(int64(request.StartTime), 0).UTC(),
		EndTime:           time.Unix(int64(request.EndTime), 0).UTC(),
		SoftCap:           amount.StringU(uint64(request.SoftCap)),
		HardCap:           amount.StringU(uint64(request.HardCap)),
		Details:           details,
		QuoteAssets:       quoteAssets,
		SaleType:          saleType,
	}
}

func getLimitsUpdateRequest(request *xdr.LimitsUpdateRequest) history.LimitsUpdateRequest {
	return history.LimitsUpdateRequest{
		DocumentHash: hex.EncodeToString(request.DocumentHash[:]),
	}
}

func getUpdateKYCRequest(request *xdr.UpdateKycRequest) history.UpdateKYCRequest {
	var kycData map[string]interface{}
	// error is ignored on purpose, we should not block ingest in case of such error
	_ = json.Unmarshal([]byte(request.KycData), &kycData)

	var externalDetails []map[string]interface{}
	for _, item := range request.ExternalDetails {
		var comment map[string]interface{}
		_ = json.Unmarshal([]byte(item), &comment)
		externalDetails = append(externalDetails, comment)
	}

	return history.UpdateKYCRequest{
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

func getReviewableRequestDetails(body *xdr.ReviewableRequestEntryBody) ([]byte, error) {
	var rawDetails interface{}
	var err error
	switch body.Type {
	case xdr.ReviewableRequestTypeAssetCreate:
		rawDetails = getAssetCreation(body.AssetCreationRequest)
	case xdr.ReviewableRequestTypeAssetUpdate:
		rawDetails = getAssetUpdate(body.AssetUpdateRequest)
	case xdr.ReviewableRequestTypeIssuanceCreate:
		rawDetails = getIssuanceRequest(body.IssuanceRequest)
	case xdr.ReviewableRequestTypePreIssuanceCreate:
		rawDetails, err = getPreIssuanceRequest(body.PreIssuanceRequest)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get pre issuance request")
		}
	case xdr.ReviewableRequestTypeWithdraw:
		rawDetails = getWithdrawalRequest(body.WithdrawalRequest)
	case xdr.ReviewableRequestTypeSale:
		rawDetails = getSaleRequest(body.SaleCreationRequest)
	case xdr.ReviewableRequestTypeLimitsUpdate:
		rawDetails = getLimitsUpdateRequest(body.LimitsUpdateRequest)
	case xdr.ReviewableRequestTypeTwoStepWithdrawal:
		rawDetails = getWithdrawalRequest(body.TwoStepWithdrawalRequest)
	case xdr.ReviewableRequestTypeAmlAlert:
		rawDetails = getAmlAlertRequest(body.AmlAlertRequest)
	case xdr.ReviewableRequestTypeUpdateKyc:
		rawDetails = getUpdateKYCRequest(body.UpdateKycRequest)
	default:
		return nil, errors.From(errors.New("unexpected reviewable request type"), map[string]interface{}{
			"request_type": body.Type.String(),
		})
	}

	details, err := json.Marshal(rawDetails)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal reviewable request details")
	}

	return details, nil
}

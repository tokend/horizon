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

	err = is.Ingestion.HistoryQ.ReviewableRequests().Insert(*histReviewableRequest)
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

	err = is.Ingestion.HistoryQ.ReviewableRequests().Update(*histReviewableRequest)
	if err != nil {
		return errors.Wrap(err, "failed to update reviewable request")
	}

	return nil
}

func convertReviewableRequest(request *xdr.ReviewableRequestEntry, ledgerCloseTime time.Time) (*history.ReviewableRequest, error) {
	var reference *string
	if request.Reference != nil {
		reference = new(string)
		*reference = string(*request.Reference)
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
		Asset:                string(request.Code),
		Policies:             int32(request.Policies),
		PreIssuedAssetSigner: request.PreissuedAssetSigner.Address(),
		MaxIssuanceAmount:    amount.StringU(uint64(request.MaxIssuanceAmount)),
		Details:              details,
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
	return history.IssuanceRequest{
		Asset:    string(request.Asset),
		Amount:   amount.StringU(uint64(request.Amount)),
		Receiver: request.Receiver.AsString(),
	}
}

func getWithdrawalRequest(request *xdr.WithdrawalRequest) history.WithdrawalRequest {
	return history.WithdrawalRequest{
		BalanceID:       request.Balance.AsString(),
		Amount:          amount.StringU(uint64(request.Amount)),
		FixedFee:        amount.StringU(uint64(request.Fee.Fixed)),
		PercentFee:      amount.StringU(uint64(request.Fee.Percent)),
		ExternalDetails: string(request.ExternalDetails),
		DestAssetCode:   string(request.Details.AutoConversion.DestAsset),
		DestAssetAmount: amount.StringU(uint64(request.Details.AutoConversion.ExpectedAmount)),
	}
}

func getSaleRequest(request *xdr.SaleCreationRequest) history.SaleRequest {
	var details map[string]interface{}
	// error is ignored on purpose, we should not block ingest in case of such error
	_ = json.Unmarshal([]byte(request.Details), &details)
	return history.SaleRequest{
		BaseAsset:  string(request.BaseAsset),
		QuoteAsset: string(request.QuoteAsset),
		StartTime:  time.Unix(int64(request.StartTime), 0).UTC(),
		EndTime:    time.Unix(int64(request.EndTime), 0).UTC(),
		Price:      amount.StringU(uint64(request.Price)),
		SoftCap:    amount.StringU(uint64(request.SoftCap)),
		HardCap:    amount.StringU(uint64(request.HardCap)),
		Details:    details,
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

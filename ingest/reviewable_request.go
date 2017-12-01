package ingest

import (
	"encoding/hex"
	"encoding/json"
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

	histReviewableRequest, err := convertReviewableRequest(reviewableRequest)
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

	histReviewableRequest, err := convertReviewableRequest(reviewableRequest)
	if err != nil {
		return errors.Wrap(err, "failed to convert reviewable request")
	}

	err = is.Ingestion.HistoryQ.ReviewableRequests().Update(*histReviewableRequest)
	if err != nil {
		return errors.Wrap(err, "failed to update reviewable request")
	}

	return nil
}

func convertReviewableRequest(request *xdr.ReviewableRequestEntry) (*history.ReviewableRequest, error) {
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
	}, nil
}

func getAssetCreation(request *xdr.AssetCreationRequest) history.AssetCreationRequest {
	return history.AssetCreationRequest{
		Asset:                 string(request.Code),
		Description:          string(request.Description),
		ExternalResourceLink: string(request.ExternalResourceLink),
		Policies:             int32(request.Policies),
		Name:                 string(request.Name),
		PreIssuedAssetSigner: request.PreissuedAssetSigner.Address(),
		MaxIssuanceAmount:    amount.StringU(uint64(request.MaxIssuanceAmount)),
	}
}

func getAssetUpdate(request *xdr.AssetUpdateRequest) history.AssetUpdateRequest {
	return history.AssetUpdateRequest{
		Asset:                 string(request.Code),
		Description:          string(request.Description),
		ExternalResourceLink: string(request.ExternalResourceLink),
		Policies:             int32(request.Policies),
	}
}

func getPreIssuanceRequest(request *xdr.PreIssuanceRequest) (history.PreIssuanceRequest, error) {
	signature, err := xdr.MarshalBase64(request.Amount)
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

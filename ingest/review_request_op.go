package ingest

import (
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/horizon/db2/history"
	"encoding/json"
)

func (is *Session) processReviewRequest(op xdr.ReviewRequestOp) {
	if is.Err != nil {
		return
	}

	var err error
	switch op.Action {
	case xdr.ReviewRequestOpActionApprove:
		err = is.approveReviewableRequest (op)
	case xdr.ReviewRequestOpActionPermanentReject:
		err = is.permanentReject(op)
	case xdr.ReviewRequestOpActionReject:
		return
	default:
		err = errors.From(errors.New("Unexpected review request action"), map[string]interface{}{
			"action_type": op.Action,
		})
	}

	if err != nil {
		is.Err = errors.Wrap(err, "failed to process review request", map[string]interface{}{
			"request_id": uint64(op.RequestId),
		})
	}
}

func (is *Session) approveReviewableRequest(op xdr.ReviewRequestOp) error {
	err := is.Cursor.HistoryQ().ReviewableRequests().Approve(uint64(op.RequestId))
	if err != nil {
		return errors.Wrap(err, "failed to approve reviewable request")
	}

	err = is.Ingestion.UpdatePayment(op.RequestId, true, nil)
	if err != nil {
		return errors.Wrap(err, "failed to approve operation")
	}

	switch op.RequestDetails.RequestType {
	case xdr.ReviewableRequestTypeWithdraw:
		err = is.setWithdrawalDetails(uint64(op.RequestId), op.RequestDetails.Withdrawal)
	}

	if err != nil {
		return errors.Wrap(err, "failed to set reviewer details")
	}

	return nil
}

func (is *Session) setWithdrawalDetails(requestID uint64, details *xdr.WithdrawalDetails) error {
	fields := logan.Field("request_id", requestID)
	request, err := is.Ingestion.HistoryQ.ReviewableRequests().ByID(requestID)
	if err != nil {
		return errors.Wrap(err, "failed to load reviewable request by id", fields)
	}

	if request == nil {
		return errors.From(errors.New("reviewable request not found"), fields)
	}

	if request.RequestType != xdr.ReviewableRequestTypeWithdraw {
		return errors.From(errors.New("expected withdrawal request"), fields.Add("request_type", request.RequestType))
	}

	var withdrawalDetails history.WithdrawalRequest
	err = json.Unmarshal(request.Details, &withdrawalDetails)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal withdrawal details", fields)
	}

	var reviewerDetails map[string]interface{}

	err = json.Unmarshal([]byte(details.ExternalDetails), &reviewerDetails)
	if err != nil {
		// we ignore here error on purpose, as it's too late to valid the error
		err = errors.Wrap(err, "failed to marshal reviewer details", fields)
		is.log.WithError(err).Warn("Reviewer sent invalid json in withdrawal details")
	}

	withdrawalDetails.ReviewerDetails = reviewerDetails
	request.Details, err = json.Marshal(withdrawalDetails)
	if err != nil {
		return errors.Wrap(err, "failed to marhsal withdrawal details", fields)
	}

	err = is.Ingestion.HistoryQ.ReviewableRequests().Update(*request)
	if err != nil {
		return errors.Wrap(err, "failed to update withdrawal request", fields)
	}

	return nil
}


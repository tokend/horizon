package ingest

import (
	"encoding/json"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/utf8"
)

func (is *Session) processReviewRequest(op xdr.ReviewRequestOp, changes xdr.LedgerEntryChanges) {
	if is.Err != nil {
		return
	}

	var err error
	switch op.Action {
	case xdr.ReviewRequestOpActionApprove:
		err = is.approveReviewableRequest(op, changes)
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

func hasKYCRequestDeleted(changes xdr.LedgerEntryChanges) bool {
	for i := range changes {
		if changes[i].Removed == nil {
			continue
		}

		if changes[i].Removed.ReviewableRequest != nil {
			return true
		}
	}

	return false
}

func (is *Session) approveReviewableRequest(op xdr.ReviewRequestOp, changes xdr.LedgerEntryChanges) error {
	// approval of two step withdrawal leads to update of request to withdrawal
	if op.RequestDetails.RequestType == xdr.ReviewableRequestTypeTwoStepWithdrawal {
		return nil
	}

	if op.RequestDetails.RequestType == xdr.ReviewableRequestTypeUpdateKyc && !hasKYCRequestDeleted(changes) {
		return nil
	}

	err := is.Ingestion.HistoryQ().ReviewableRequests().Approve(uint64(op.RequestId))
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
	request, err := is.Ingestion.HistoryQ().ReviewableRequests().ByID(requestID)
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

	externalDetails := utf8.Scrub(details.ExternalDetails)
	err = json.Unmarshal([]byte(externalDetails), &reviewerDetails)
	if err != nil {
		// we ignore here error on purpose, as it's too late to valid the error
		err = errors.Wrap(err, "failed to marshal reviewer details", fields)
		is.log.WithError(err).WithFields(logan.F{
			"scrubbed_details": externalDetails,
			"original_details": details.ExternalDetails,
		}).Warn("Reviewer sent invalid json in withdrawal details")
	}

	withdrawalDetails.ReviewerDetails = reviewerDetails
	request.Details, err = json.Marshal(withdrawalDetails)
	if err != nil {
		return errors.Wrap(err, "failed to marhsal withdrawal details", fields)
	}

	err = is.Ingestion.HistoryQ().ReviewableRequests().Update(*request)
	if err != nil {
		return errors.Wrap(err, "failed to update withdrawal request", fields)
	}

	return nil
}

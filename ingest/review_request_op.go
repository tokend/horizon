package ingest

import (
	"encoding/json"

	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/utf8"
	"gitlab.com/tokend/go/xdr"
)

func (is *Session) processReviewRequest(op xdr.ReviewRequestOp, changes xdr.LedgerEntryChanges) (err error) {
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
		return errors.Wrap(err, "failed to process review request", map[string]interface{}{
			"request_id": uint64(op.RequestId),
		})
	}
	return nil
}

func hasDeletedReviewableRequest(changes xdr.LedgerEntryChanges) bool {
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

func (is *Session) approveReviewableRequest(op xdr.ReviewRequestOp, changes xdr.LedgerEntryChanges) (err error) {
	// approval of two step withdrawal leads to update of request to withdrawal
	if op.RequestDetails.RequestType == xdr.ReviewableRequestTypeTwoStepWithdrawal {
		return nil
	}

	if op.RequestDetails.RequestType == xdr.ReviewableRequestTypeUpdateKyc && !hasDeletedReviewableRequest(changes) {
		return nil
	}

	err = is.Ingestion.HistoryQ().ReviewableRequests().Approve(uint64(op.RequestId))
	if err != nil {
		return errors.Wrap(err, "failed to approve reviewable request")
	}

	switch op.RequestDetails.RequestType {
	case xdr.ReviewableRequestTypeWithdraw:
		err = is.setWithdrawalDetails(uint64(op.RequestId), op.RequestDetails.Withdrawal)
	case xdr.ReviewableRequestTypeContract:
		err = is.processContractLedgerChanges(nil, nil)
	case xdr.ReviewableRequestTypeInvoice:
		err = is.setWaitingForConfirmationState(uint64(op.RequestId))
	}

	if err != nil {
		return errors.Wrap(err, "failed to set request details")
	}

	return nil
}

func (is *Session) setWithdrawalDetails(requestID uint64, details *xdr.WithdrawalDetails) (err error) {
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

	var withdrawalDetails *history.WithdrawalRequest
	if request.Details.Withdraw != nil {
		withdrawalDetails = request.Details.Withdraw
	} else if request.Details.TwoStepWithdraw != nil {
		withdrawalDetails = request.Details.TwoStepWithdraw
	} else {
		return errors.New("Unexpected state: expected withdrawal details to be available")
	}

	withdrawalDetails.ReviewerDetails = reviewerDetails
	err = is.Ingestion.HistoryQ().ReviewableRequests().Update(*request)
	if err != nil {
		return errors.Wrap(err, "failed to update withdrawal request", fields)
	}

	return nil
}

func (is *Session) setWaitingForConfirmationState(requestID uint64) error {
	request, err := is.Ingestion.HistoryQ().ReviewableRequests().ByID(requestID)
	if err != nil {
		return errors.Wrap(err, "failed to get request", logan.F{
			"request_id": requestID,
		})
	}

	if (request == nil) || (request.Details.Invoice == nil) || (request.Details.Invoice.ContractID == nil) {
		return nil
	}

	err = is.Ingestion.HistoryQ().ReviewableRequests().UpdateStates([]int64{int64(requestID)},
		history.ReviewableRequestStateWaitingForConfirmation)
	if err != nil {
		return errors.Wrap(err, "failed to update request state")
	}

	return nil
}

func (is *Session) processContractLedgerChanges(isDisputeStart, isRevert *bool) error {
	ledgerChanges := is.Cursor.OperationChanges()
	for _, change := range ledgerChanges {
		switch change.Type {
		case xdr.LedgerEntryChangeTypeCreated:
			if change.Created.Data.Type != xdr.LedgerEntryTypeContract {
				continue
			}

			contract := convertContract(change.Created.Data.MustContract())

			err := is.Ingestion.Contracts(contract)
			if err != nil {
				return errors.Wrap(err, "failed to ingest contract", logan.F{
					"contract_id": uint64(change.Created.Data.MustContract().ContractId),
				})
			}
		case xdr.LedgerEntryChangeTypeUpdated:
			if change.Updated.Data.Type != xdr.LedgerEntryTypeContract {
				continue
			}

			contract := convertContract(change.Updated.Data.MustContract())

			err := is.Ingestion.HistoryQ().Contracts().Update(contract)
			if err != nil {
				return errors.Wrap(err, "failed to update contract", logan.F{
					"contract_id": uint64(change.Updated.Data.MustContract().ContractId),
				})
			}
		case xdr.LedgerEntryChangeTypeRemoved:
			if change.Removed.Type != xdr.LedgerEntryTypeContract {
				continue
			}

			contractID := int64(change.Removed.Contract.ContractId)

			contract, err := is.Ingestion.HistoryQ().Contracts().ByID(contractID)
			if err != nil {
				return errors.Wrap(err, "failed to get contract", logan.F{
					"contract_id": contractID,
				})
			}

			if isRevert == nil {
				return is.updateInvoicesContractStates(contract,
					history.ReviewableRequestStateApproved,
					int32(xdr.ContractStateCustomerConfirmed)|
						int32(xdr.ContractStateContractorConfirmed),
				)
			}

			if *isRevert {
				return is.updateInvoicesContractStates(contract,
					history.ReviewableRequestStatePermanentlyRejected,
					int32(xdr.ContractStateRevertingResolve),
				)
			}

			contract.State |= int32(xdr.ContractStateNotRevertingResolve)
			err = is.Ingestion.HistoryQ().Contracts().Update(contract)
			if err != nil {
				return errors.Wrap(err, "failed to update contract state")
			}

			return is.updateNotRevertingContractInovices(contract)
		}
	}
	return nil
}

func (is *Session) updateInvoicesContractStates(
	contract history.Contract,
	invoiceState history.ReviewableRequestState,
	contactState int32,
) error {
	err := is.Ingestion.HistoryQ().ReviewableRequests().UpdateStates(
		contract.Invoices,
		invoiceState,
	)
	if err != nil {
		return errors.Wrap(err, "failed to update invoices states")
	}

	contract.State |= contactState
	err = is.Ingestion.HistoryQ().Contracts().Update(contract)
	if err != nil {
		return errors.Wrap(err, "failed to update contract state")
	}

	return nil
}

func (is *Session) updateNotRevertingContractInovices(contract history.Contract) error {
	invoices, err := is.Ingestion.HistoryQ().ReviewableRequests().ByIDs(contract.Invoices).Select()
	if err != nil {
		return errors.Wrap(err, "failed to load invoices")
	}

	var approvedInvoicesIDs, rejectedInvoicesIDs []int64

	for _, invoice := range invoices {
		if invoice.RequestState == history.ReviewableRequestStateWaitingForConfirmation {
			approvedInvoicesIDs = append(approvedInvoicesIDs, invoice.ID)
			continue
		}

		rejectedInvoicesIDs = append(rejectedInvoicesIDs, invoice.ID)
	}

	err = is.Ingestion.HistoryQ().ReviewableRequests().UpdateStates(approvedInvoicesIDs,
		history.ReviewableRequestStateApproved)
	if err != nil {
		return errors.Wrap(err, "failed to update approved invoices", logan.F{
			"request_ids": approvedInvoicesIDs,
		})
	}

	err = is.Ingestion.HistoryQ().ReviewableRequests().UpdateStates(rejectedInvoicesIDs,
		history.ReviewableRequestStatePermanentlyRejected)
	if err != nil {
		return errors.Wrap(err, "failed to update rejected invoices", logan.F{
			"request_ids": rejectedInvoicesIDs,
		})
	}

	return nil
}

func convertContract(rawContract xdr.ContractEntry) history.Contract {
	disputer := ""
	var disputeReason map[string]interface{}
	if (int32(rawContract.State) & int32(xdr.ContractStateDisputing)) != 0 {
		disputer = rawContract.DisputeDetails.Disputer.Address()
		// error is ignored on purpose, we should not block ingest in case of such error
		_ = json.Unmarshal([]byte(rawContract.DisputeDetails.Reason), &disputeReason)
	}

	var details []db2.Details
	for _, item := range rawContract.Details {
		var detail map[string]interface{}
		_ = json.Unmarshal([]byte(item), &detail)
		details = append(details, detail)
	}

	var invoices []int64
	for _, item := range rawContract.InvoiceRequestsIDs {
		invoices = append(invoices, int64(item))
	}

	return history.Contract{
		TotalOrderID: db2.TotalOrderID{
			ID: int64(rawContract.ContractId),
		},
		Contractor:    rawContract.Contractor.Address(),
		Customer:      rawContract.Customer.Address(),
		Escrow:        rawContract.Escrow.Address(),
		Disputer:      disputer,
		StartTime:     time.Unix(int64(rawContract.StartTime), 0).UTC(),
		EndTime:       time.Unix(int64(rawContract.EndTime), 0).UTC(),
		Details:       details,
		Invoices:      invoices,
		DisputeReason: disputeReason,
		State:         int32(rawContract.State),
	}
}

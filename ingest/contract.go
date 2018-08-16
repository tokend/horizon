package ingest

import (
	"encoding/json"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/go/xdr"
)

func contractCreate(is *Session, entry *xdr.LedgerEntry) error {
	contract := convertContract(entry.Data.MustContract())

	err := is.Ingestion.Contracts(contract)
	if err != nil {
		return errors.Wrap(err, "failed to ingest contract", logan.F{
			"contract_id": uint64(entry.Data.MustContract().ContractId),
		})
	}

	return nil
}

func contractUpdate(is *Session, entry *xdr.LedgerEntry) error {
	contract := convertContract(entry.Data.MustContract())

	err := is.Ingestion.HistoryQ().Contracts().Update(contract)
	if err != nil {
		return errors.Wrap(err, "failed to update contract", logan.F{
			"contract_id": uint64(entry.Data.MustContract().ContractId),
		})
	}

	return nil
}

func contractDelete(is *Session, key *xdr.LedgerKey) error {
	manageContractOp := is.Cursor.Operation().Body.MustManageContractOp()
	if manageContractOp.ContractId != key.Contract.ContractId {
		return errors.New("Expected manage contract op contract id to be equal ledger key contract id")
	}
	isRevert := manageContractOp.Data.IsRevert
	contractID := int64(manageContractOp.ContractId)

	if isRevert == nil {
		return is.updateInvoicesContractStates(contractID,
			history.ReviewableRequestStateApproved,
			int32(xdr.ContractStateCustomerConfirmed)|
				int32(xdr.ContractStateContractorConfirmed),
			[]history.ReviewableRequestState{history.ReviewableRequestStateApproved})
	}

	if *isRevert {
		return is.updateInvoicesContractStates(contractID,
			history.ReviewableRequestStatePermanentlyRejected,
			int32(xdr.ContractStateRevertingResolve),
			[]history.ReviewableRequestState{history.ReviewableRequestStatePending,
				history.ReviewableRequestStateWaitingForConfirmation})
	}

	err := is.Ingestion.HistoryQ().Contracts().AddState(contractID, int32(xdr.ContractStateNotRevertingResolve))
	if err != nil {
		return errors.Wrap(err, "failed to update contract state")
	}

	return is.updateNotRevertingContractInvoices(contractID)
}

func (is *Session) updateInvoicesContractStates(
	contractID int64,
	invoiceState history.ReviewableRequestState,
	contractState int32,
	oldStates []history.ReviewableRequestState,
) error {
	err := is.Ingestion.HistoryQ().ReviewableRequests().
		UpdateInvoicesStates(invoiceState, contractID, oldStates)
	if err != nil {
		return errors.Wrap(err, "failed to update invoices states")
	}

	err = is.Ingestion.HistoryQ().Contracts().AddState(contractID, contractState)
	if err != nil {
		return errors.Wrap(err, "failed to update contract state")
	}

	return nil
}

func (is *Session) updateNotRevertingContractInvoices(contractID int64) error {
	err := is.Ingestion.HistoryQ().ReviewableRequests().UpdateInvoicesStates(
		history.ReviewableRequestStateApproved, contractID,
		[]history.ReviewableRequestState{history.ReviewableRequestStateWaitingForConfirmation})
	if err != nil {
		return errors.Wrap(err, "failed to update approved invoices", logan.F{
			"contract_id": contractID,
		})
	}

	err = is.Ingestion.HistoryQ().ReviewableRequests().UpdateInvoicesStates(
		history.ReviewableRequestStatePermanentlyRejected, contractID,
		[]history.ReviewableRequestState{history.ReviewableRequestStatePending})
	if err != nil {
		return errors.Wrap(err, "failed to update rejected invoices", logan.F{
			"contract_id": contractID,
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

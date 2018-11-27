package changes

import (
	"encoding/json"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/history"
)

type contractStorage interface {
	//inserts contract into DB
	InsertContract(contract history.Contract) error
	//updates contract
	UpdateContract(contract history.Contract) error
	//Adds state to contract
	AddContractState(id uint64, state uint64) error
}

type contractChanges struct {
	storage    contractStorage
	reqStorage reviewableRequestStorage
}

func (c *contractChanges) Created(lc LedgerChange) error {
	rawContract := lc.LedgerChange.MustCreated().Data.MustContract()
	contract := c.convertContract(rawContract)

	err := c.storage.InsertContract(contract)
	if err != nil {
		return errors.Wrap(err, "failed to insert contract", logan.F{
			"contract_id": contract.ID,
		})
	}

	return nil
}

func (c *contractChanges) Updated(lc LedgerChange) error {
	rawContract := lc.LedgerChange.MustUpdated().Data.MustContract()
	contract := c.convertContract(rawContract)

	err := c.storage.UpdateContract(contract)
	if err != nil {
		return errors.Wrap(err, "failed to update contract", logan.F{
			"contract_id": contract.ID,
		})
	}

	return nil
}

func (c *contractChanges) Deleted(lc LedgerChange) error {
	contractKey := lc.LedgerChange.MustRemoved().MustContract()
	contractID := contractKey.ContractId
	manageContractOp := lc.Operation.Body.MustManageContractOp()
	if contractID != manageContractOp.ContractId {
		return errors.New("Expected ledger key and manage contract op to have same contract id")
	}

	switch manageContractOp.Data.Action {
	case xdr.ManageContractActionConfirmCompleted:
		err := c.storage.AddContractState(uint64(contractID), uint64(xdr.ContractStateCustomerConfirmed|xdr.ContractStateContractorConfirmed))
		if err != nil {
			return errors.Wrap(err, "failed to update invoices", logan.F{
				"contract_id": contractID,
			})
		}
		return nil

	case xdr.ManageContractActionResolveDispute:
		isRevert := manageContractOp.Data.IsRevert
		if isRevert != nil && *isRevert {
			err := c.processRevert(uint64(contractID))
			if err != nil {
				return errors.Wrap(err, "failed to process contract revert")
			}
			return nil
		}

		err := c.processResolveDispute(uint64(contractID))

		if err != nil {
			return errors.Wrap(err, "failed to process contract resolve dispute")
		}

		return nil
	default:
		return errors.New("Unexpected case in delete contract")
	}
}

func (c *contractChanges) processResolveDispute(id uint64) error {
	err := c.reqStorage.UpdateInvoices(
		id,
		uint64(history.ReviewableRequestStateWaitingForConfirmation),
		uint64(history.ReviewableRequestStateApproved),
	)

	if err != nil {
		return errors.Wrap(err, "failed to update approved invoices", logan.F{
			"contract_id": id,
		})
	}

	err = c.reqStorage.UpdateInvoices(
		id,
		uint64(history.ReviewableRequestStatePending),
		uint64(history.ReviewableRequestStatePermanentlyRejected),
	)
	if err != nil {
		return errors.Wrap(err, "failed to update rejected invoices", logan.F{
			"contract_id": id,
		})
	}
	return nil
}

func (c *contractChanges) processRevert(id uint64) error {
	err := c.reqStorage.UpdateInvoices(
		id,
		uint64(history.ReviewableRequestStateWaitingForConfirmation|
			history.ReviewableRequestStatePending),
		uint64(history.ReviewableRequestStatePermanentlyRejected),
	)

	if err != nil {
		return errors.Wrap(err, "failed to update invoices", logan.F{
			"contract_id": id,
		})
	}

	err = c.storage.AddContractState(id, uint64(xdr.ContractStateRevertingResolve))
	if err != nil {
		return errors.Wrap(err, "failed to add contract state", logan.F{
			"contract_id": id,
		})
	}

	return nil
}

func (c *contractChanges) convertContract(rawContract xdr.ContractEntry) history.Contract {
	var initialDetails map[string]interface{}
	_ = json.Unmarshal([]byte(string(rawContract.InitialDetails)), &initialDetails)

	var customerDetails map[string]interface{}
	if rawContract.Ext.V == xdr.LedgerVersionAddCustomerDetailsToContract {
		_ = json.Unmarshal([]byte(string(rawContract.Ext.MustCustomerDetails())), &customerDetails)
	}

	var invoices []int64
	for _, item := range rawContract.InvoiceRequestsIDs {
		invoices = append(invoices, int64(item))
	}

	return history.Contract{
		TotalOrderID: db2.TotalOrderID{
			ID: int64(rawContract.ContractId),
		},
		Contractor:      rawContract.Contractor.Address(),
		Customer:        rawContract.Customer.Address(),
		Escrow:          rawContract.Escrow.Address(),
		StartTime:       unixToTime(int64(rawContract.StartTime)),
		EndTime:         unixToTime(int64(rawContract.EndTime)),
		InitialDetails:  initialDetails,
		CustomerDetails: customerDetails,
		Invoices:        invoices,
		State:           int32(rawContract.State),
	}
}

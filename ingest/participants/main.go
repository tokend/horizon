// Package participants contains functions to derive a set of "participant"
// addresses for various data structures in the Stellar network's ledger.
package participants

import (
	"fmt"

	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/go/xdr"
)

// ForOperation returns all the participating accounts from the
// provided operation.

type Participant struct {
	AccountID xdr.AccountId
	BalanceID *xdr.BalanceId
	Details   interface{}
}

func ForOperation(
	DB *db2.Repo,
	tx *xdr.Transaction,
	op *xdr.Operation,
	opResult xdr.OperationResultTr,
	ledgerChanges xdr.LedgerEntryChanges,
	ledger *core.LedgerHeader,
) (result []Participant, err error) {
	sourceParticipant := &Participant{}
	if op.SourceAccount != nil {
		sourceParticipant.AccountID = *op.SourceAccount
	} else {
		sourceParticipant.AccountID = tx.SourceAccount
	}
	switch op.Body.Type {
	case xdr.OperationTypeCreateAccount:
		result = append(result, Participant{op.Body.MustCreateAccountOp().Destination, nil, nil})
	case xdr.OperationTypePayment:
		paymentOp := op.Body.MustPaymentOp()
		paymentResponse := opResult.MustPaymentResult().MustPaymentResponse()

		if paymentOp.InvoiceReference != nil {
			sourceParticipant = nil
			break
		}

		result = append(result, Participant{paymentResponse.Destination, &paymentOp.DestinationBalanceId, nil})
		sourceParticipant.BalanceID = &paymentOp.SourceBalanceId
	case xdr.OperationTypeSetOptions:
	// the only direct participant is the source_account
	case xdr.OperationTypeSetFees:
	// the only direct participant is the source_account
	case xdr.OperationTypeManageAccount:
		manageAccountOp := op.Body.MustManageAccountOp()
		result = append(result, Participant{manageAccountOp.Account, nil, nil})
	case xdr.OperationTypeCreateWithdrawalRequest:
		createWithdrawalRequest := op.Body.MustCreateWithdrawalRequestOp()
		sourceParticipant.BalanceID = &createWithdrawalRequest.Request.Balance
	case xdr.OperationTypeManageBalance:
		manageBalanceOp := op.Body.MustManageBalanceOp()
		if sourceParticipant.AccountID.Address() != manageBalanceOp.Destination.Address() {
			result = append(result, Participant{manageBalanceOp.Destination, nil, nil})
		}
	case xdr.OperationTypeReviewPaymentRequest:
	// the only direct participant is the source_account
	case xdr.OperationTypeManageAsset:
	// the only direct participant is the source_accountWWW
	case xdr.OperationTypeSetLimits:
		setLimitsOp := op.Body.MustSetLimitsOp()
		if setLimitsOp.Account != nil {
			details := map[string]interface{}{}
			details["account"] = setLimitsOp.Account
		}
	case xdr.OperationTypeDirectDebit:
		debitOp := op.Body.MustDirectDebitOp()
		paymentOp := debitOp.PaymentOp
		paymentResponse := opResult.MustDirectDebitResult().MustSuccess().PaymentResponse
		details := map[string]interface{}{}
		details["initiated_by"] = sourceParticipant.AccountID
		result = append(result, Participant{paymentResponse.Destination, &paymentOp.DestinationBalanceId, &details})
		sourceParticipant.BalanceID = &paymentOp.SourceBalanceId
		sourceParticipant.AccountID = debitOp.From
	case xdr.OperationTypeManageAssetPair:
		// the only direct participant is the source_account
	case xdr.OperationTypeManageOffer:
		manageOfferOp := op.Body.MustManageOfferOp()
		manageOfferResult := opResult.MustManageOfferResult()
		result = addMatchParticipants(result, sourceParticipant.AccountID, manageOfferOp.BaseBalance,
			manageOfferOp.QuoteBalance, manageOfferOp.IsBuy, manageOfferResult.Success)
		sourceParticipant = nil
	case xdr.OperationTypeManageInvoice:
		manageInvoiceOp := op.Body.MustManageInvoiceOp()
		if manageInvoiceOp.Amount == 0 {
			sourceParticipant = nil
			break
		}
		sourceParticipant.BalanceID = &manageInvoiceOp.ReceiverBalance
		result = append(result, Participant{manageInvoiceOp.Sender, &opResult.ManageInvoiceResult.Success.SenderBalance, nil})
	case xdr.OperationTypeReviewRequest:
		request := getReviewableRequestByID(uint64(op.Body.MustReviewRequestOp().RequestId), ledgerChanges)
		if request != nil && sourceParticipant.AccountID.Address() != request.Requestor.Address() {
			result = append(result, Participant{
				AccountID: request.Requestor,
				BalanceID: nil,
				Details:   nil,
			})
		}
	case xdr.OperationTypeCreatePreissuanceRequest:
		// the only direct participant is the source_account
	case xdr.OperationTypeCreateIssuanceRequest:
		manageIssuanceRequest := op.Body.MustCreateIssuanceRequestOp()
		manageIssuanceResult := opResult.MustCreateIssuanceRequestResult()
		result = append(result, Participant{manageIssuanceResult.MustSuccess().Receiver,
			&manageIssuanceRequest.Request.Receiver, nil})
	case xdr.OperationTypeCreateSaleRequest:
		// the only direct participant is the source_account
	case xdr.OperationTypeCheckSaleState:
		manageOfferResult := opResult.MustCheckSaleStateResult()
		saleClosed := manageOfferResult.Success.Effect.SaleClosed
		if saleClosed == nil {
			break
		}

		for i := range saleClosed.Results {
			result = addMatchParticipants(result, saleClosed.SaleOwner, saleClosed.Results[i].SaleBaseBalance,
				saleClosed.Results[i].SaleQuoteBalance, false, &saleClosed.Results[i].SaleDetails)
		}

		sourceParticipant = nil
	case xdr.OperationTypeManageExternalSystemAccountIdPoolEntry:
		// the only direct participant is the source_account
	case xdr.OperationTypeBindExternalSystemAccountId:
		// the only direct participant is the source_account
	case xdr.OperationTypeCreateAmlAlert:
		// TODO add participant
	case xdr.OperationTypeCreateKycRequest:
		updateKYCRequestData := op.Body.MustCreateUpdateKycRequestOp().UpdateKycRequestData
		if sourceParticipant.AccountID.Address() != updateKYCRequestData.AccountToUpdateKyc.Address() {
			result = append(result, Participant{
				AccountID: updateKYCRequestData.AccountToUpdateKyc,
				BalanceID: nil,
				Details:   nil,
			})
		}
	default:
		err = fmt.Errorf("unknown operation type: %s", op.Body.Type)
	}

	if sourceParticipant != nil {
		result = append(result, *sourceParticipant)
	}
	return
}

func getReviewableRequestByID(id uint64, ledgerChanges xdr.LedgerEntryChanges) *xdr.ReviewableRequestEntry {
	for i := range ledgerChanges {
		ledgerChange := ledgerChanges[i]
		var reviewableRequest *xdr.ReviewableRequestEntry
		switch ledgerChange.Type {
		case xdr.LedgerEntryChangeTypeCreated:
			reviewableRequest = ledgerChange.Created.Data.ReviewableRequest
		case xdr.LedgerEntryChangeTypeUpdated:
			reviewableRequest = ledgerChange.Updated.Data.ReviewableRequest
		default:
			continue
		}

		if reviewableRequest == nil {
			continue
		}

		if uint64(reviewableRequest.RequestId) == id {
			return reviewableRequest
		}
	}

	return nil
}

func addMatchParticipants(participants []Participant, offerSourceID xdr.AccountId, baseBalanceID xdr.BalanceId,
	quoteBalanceID xdr.BalanceId, isBuy bool, result *xdr.ManageOfferSuccessResult) []Participant {
	if result == nil || len(result.OffersClaimed) == 0 {
		return participants
	}

	matchesByBalance := NewMatchesDetailsByBalance()

	for _, offerClaimed := range result.OffersClaimed {

		claimedOfferMatch := NewMatch(offerClaimed.BaseAmount, offerClaimed.QuoteAmount, offerClaimed.BFeePaid, offerClaimed.CurrentPrice)
		matchesByBalance.Add(offerClaimed.BAccountId, []xdr.BalanceId{offerClaimed.BaseBalance, offerClaimed.QuoteBalance},
			result.BaseAsset, result.QuoteAsset, !isBuy, claimedOfferMatch)

		offerMatch := NewMatch(offerClaimed.BaseAmount, offerClaimed.QuoteAmount, offerClaimed.AFeePaid, offerClaimed.CurrentPrice)
		matchesByBalance.Add(offerSourceID, []xdr.BalanceId{baseBalanceID, quoteBalanceID}, result.BaseAsset, result.QuoteAsset, isBuy, offerMatch)

	}

	return matchesByBalance.ToParticipants(participants)
}

// ForTransaction returns all the participating accounts from the provided
// transaction.
func ForTransaction(
	DB *db2.Repo,
	tx *xdr.Transaction,
	opResults []xdr.OperationResult,
	meta *xdr.TransactionMeta,
	ledgerChanges *xdr.LedgerEntryChanges,
	ledger *core.LedgerHeader,
) (result []xdr.AccountId, err error) {

	result = append(result, tx.SourceAccount)

	p, err := forMeta(meta)
	if err != nil {
		return
	}
	result = append(result, p...)

	p, err = forChanges(ledgerChanges)
	if err != nil {
		return
	}
	result = append(result, p...)

	for i := range tx.Operations {
		participants, err := ForOperation(DB, tx, &tx.Operations[i], *opResults[i].Tr, nil, ledger)
		if err != nil {
			return nil, err
		}
		for _, participant := range participants {
			result = append(result, participant.AccountID)
		}
	}

	result = dedupe(result)
	return
}

func getAccountIDByBalance(q history.Q, balanceID string) (result *xdr.AccountId, err error) {
	var targetBalance history.Balance
	err = q.BalanceByID(&targetBalance, balanceID)
	if err != nil {
		return nil, err
	}
	var aid xdr.AccountId
	aid.SetAddress(targetBalance.AccountID)
	return &aid, nil
}

// dedupe remove any duplicate ids from `in`
func dedupe(in []xdr.AccountId) (out []xdr.AccountId) {
	set := map[string]xdr.AccountId{}
	for _, id := range in {
		set[id.Address()] = id
	}

	for _, id := range set {
		out = append(out, id)
	}
	return
}

func forChanges(
	changes *xdr.LedgerEntryChanges,
) (result []xdr.AccountId, err error) {

	for _, c := range *changes {
		var account *xdr.AccountId

		switch c.Type {
		case xdr.LedgerEntryChangeTypeCreated:
			account = forLedgerEntry(c.MustCreated())
		case xdr.LedgerEntryChangeTypeRemoved:
			account = forLedgerKey(c.MustRemoved())
		case xdr.LedgerEntryChangeTypeUpdated:
			account = forLedgerEntry(c.MustUpdated())
		case xdr.LedgerEntryChangeTypeState:
			account = forLedgerEntry(c.MustState())
		default:
			err = fmt.Errorf("Unknown change type: %s", c.Type)
			return
		}

		if account != nil {
			result = append(result, *account)
		}
	}

	return
}

func forLedgerEntry(le xdr.LedgerEntry) *xdr.AccountId {
	if le.Data.Type != xdr.LedgerEntryTypeAccount {
		return nil
	}
	aid := le.Data.MustAccount().AccountId
	return &aid
}

func forLedgerKey(lk xdr.LedgerKey) *xdr.AccountId {
	if lk.Type != xdr.LedgerEntryTypeAccount {
		return nil
	}
	aid := lk.MustAccount().AccountId
	return &aid
}

func forMeta(
	meta *xdr.TransactionMeta,
) (result []xdr.AccountId, err error) {

	if meta.Operations == nil {
		return
	}

	for _, op := range *meta.Operations {
		var acc []xdr.AccountId
		acc, err = forChanges(&op.Changes)
		if err != nil {
			return
		}

		result = append(result, acc...)
	}

	return
}

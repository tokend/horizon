// Package participants contains functions to derive a set of "participant"
// addresses for various data structures in the Stellar network's ledger.
package participants

import (
	"fmt"

	"time"

	"bullioncoin.githost.io/development/go/amount"
	"bullioncoin.githost.io/development/go/xdr"
	"gitlab.com/distributed_lab/tokend/horizon/db2"
	"gitlab.com/distributed_lab/tokend/horizon/db2/core"
	"gitlab.com/distributed_lab/tokend/horizon/db2/history"
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
	ledger *core.LedgerHeader,
) (result []Participant, err error) {
	q := history.Q{Repo: DB}
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
	case xdr.OperationTypeManageCoinsEmissionRequest:
		if !opResult.MustManageCoinsEmissionRequestResult().ManageRequestInfo.Fulfilled {
			break
		}

		balance := op.Body.MustManageCoinsEmissionRequestOp().Receiver
		accountID, err := getAccountIDByBalance(q, balance.AsString())
		if err != nil {
			return nil, err
		}
		result = append(result, Participant{*accountID, &balance, nil})
	case xdr.OperationTypeReviewCoinsEmissionRequest:
		balance := op.Body.MustReviewCoinsEmissionRequestOp().Request.Receiver
		accountID, err := getAccountIDByBalance(q, balance.AsString())
		if err != nil {
			return nil, err
		}
		result = append(result, Participant{*accountID, &balance, nil})
	case xdr.OperationTypeSetFees:
	// the only direct participant is the source_account
	case xdr.OperationTypeManageAccount:
		manageAccountOp := op.Body.MustManageAccountOp()
		result = append(result, Participant{manageAccountOp.Account, nil, nil})
	case xdr.OperationTypeForfeit:
		balance := op.Body.MustForfeitOp().Balance
		accountID, err := getAccountIDByBalance(q, balance.AsString())
		if err != nil {
			return nil, err
		}
		result = append(result, Participant{*accountID, &balance, nil})
	case xdr.OperationTypeManageForfeitRequest:
		manageForfeitRequestOp := op.Body.MustManageForfeitRequestOp()
		sourceParticipant.BalanceID = &manageForfeitRequestOp.Balance
		manageForfeitRequestResult := opResult.MustManageForfeitRequestResult()
		result = append(result, Participant{manageForfeitRequestResult.ForfeitRequestDetails.Exchange, nil, nil})
	case xdr.OperationTypeRecover:
		result = append(result, Participant{op.Body.MustRecoverOp().Account, nil, nil})
	case xdr.OperationTypeManageBalance:
		manageBalanceOp := op.Body.MustManageBalanceOp()
		result = append(result, Participant{manageBalanceOp.Destination, &manageBalanceOp.BalanceId, nil})
	case xdr.OperationTypeReviewPaymentRequest:
	// the only direct participant is the source_account
	case xdr.OperationTypeManageAsset:
	// the only direct participant is the source_account
	case xdr.OperationTypeDemurrage:
		demurrageOp := op.Body.MustDemurrageOp()
		demurrageResult := opResult.MustDemurrageResult()
		for _, balanceDemurrage := range demurrageResult.DemurrageInfo.Demurrages {
			details := map[string]interface{}{}
			details["amount"] = amount.String(int64(balanceDemurrage.Amount))
			details["asset"] = demurrageOp.Asset
			details["period_to"] = time.Unix(ledger.CloseTime, 0).UTC()
			details["period_from"] = time.Unix(ledger.CloseTime, 0).Add(time.Duration(-int64(balanceDemurrage.Period)) * time.Second).UTC()
			tmpBalance := balanceDemurrage.Balance
			result = append(result, Participant{balanceDemurrage.Account, &tmpBalance, &details})
		}
	case xdr.OperationTypeUploadPreemissions:
		// the only direct participant is the source_account
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
		result = addManageOfferParticipants(result, sourceParticipant.AccountID, manageOfferOp, manageOfferResult)
		sourceParticipant = nil
	case xdr.OperationTypeManageInvoice:
		manageInvoiceOp := op.Body.MustManageInvoiceOp()
		if manageInvoiceOp.Amount == 0 {
			sourceParticipant = nil
			break
		}
		sourceParticipant.BalanceID = &manageInvoiceOp.ReceiverBalance
		result = append(result, Participant{manageInvoiceOp.Sender, &opResult.MustManageInvoiceResult().Success.SenderBalance, nil})
	default:
		err = fmt.Errorf("Unknown operation type: %s", op.Body.Type)
	}

	if sourceParticipant != nil {
		result = append(result, *sourceParticipant)
	}
	return
}

func addManageOfferParticipants(participants []Participant, sourceID xdr.AccountId, op xdr.ManageOfferOp, result xdr.ManageOfferResult) []Participant {
	if result.Success == nil || len(result.Success.OffersClaimed) == 0 {
		return participants
	}

	matchesByBalance := NewMatchesDetailsByBalance()

	for _, offerClaimed := range result.Success.OffersClaimed {

		claimedOfferMatch := NewMatch(offerClaimed.BaseAmount, offerClaimed.QuoteAmount, offerClaimed.BFeePaid, offerClaimed.CurrentPrice)
		matchesByBalance.Add(offerClaimed.BAccountId, offerClaimed.BaseBalance, result.Success.BaseAsset, result.Success.QuoteAsset, !op.IsBuy, claimedOfferMatch)
		matchesByBalance.Add(offerClaimed.BAccountId, offerClaimed.QuoteBalance, result.Success.BaseAsset, result.Success.QuoteAsset, !op.IsBuy, claimedOfferMatch)

		offerMatch := NewMatch(offerClaimed.BaseAmount, offerClaimed.QuoteAmount, offerClaimed.AFeePaid, offerClaimed.CurrentPrice)
		matchesByBalance.Add(sourceID, op.BaseBalance, result.Success.BaseAsset, result.Success.QuoteAsset, op.IsBuy, offerMatch)
		matchesByBalance.Add(sourceID, op.QuoteBalance, result.Success.BaseAsset, result.Success.QuoteAsset, op.IsBuy, offerMatch)

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
	feeMeta *xdr.LedgerEntryChanges,
	ledger *core.LedgerHeader,
) (result []xdr.AccountId, err error) {

	result = append(result, tx.SourceAccount)

	p, err := forMeta(meta)
	if err != nil {
		return
	}
	result = append(result, p...)

	p, err = forChanges(feeMeta)
	if err != nil {
		return
	}
	result = append(result, p...)

	for i := range tx.Operations {
		participants, err := ForOperation(DB, tx, &tx.Operations[i], *opResults[i].Tr, ledger)
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
		case xdr.LedgerEntryChangeTypeLedgerEntryCreated:
			account = forLedgerEntry(c.MustCreated())
		case xdr.LedgerEntryChangeTypeLedgerEntryRemoved:
			account = forLedgerKey(c.MustRemoved())
		case xdr.LedgerEntryChangeTypeLedgerEntryUpdated:
			account = forLedgerEntry(c.MustUpdated())
		case xdr.LedgerEntryChangeTypeLedgerEntryState:
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

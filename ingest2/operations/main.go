package operations

import (
	"encoding/json"

	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type providerCluster interface {
	// GetOperationIDProvider returns operationIDProvider
	GetOperationIDProvider() operationIDProvider
	// GetParticipantEffectIDProvider returns participantEffectIDProvider
	GetParticipantEffectIDProvider() participantEffectIDProvider
	// GetPubKeyProvider returns publicKeyProvider
	GetPubKeyProvider() publicKeyProvider
	// GetBalanceProvider returns balanceProvider
	GetBalanceProvider() balanceProvider
	// GetRequestProvider returns requestProvider
	GetRequestProvider() requestProvider
}

type operationIDProvider interface {
	// GetOperationID returns unique id of current operation
	GetOperationID() int64
}

type participantEffectIDProvider interface {
	// GetNextParticipantEffectID return unique value for participant effect
	GetNextParticipantEffectID() int64
}

type publicKeyProvider interface {
	// GetAccountID returns int value which corresponds to xdr.AccountId
	GetAccountID(raw xdr.AccountId) int64
	// GetBalanceID returns int value which corresponds to xdr.BalanceId
	GetBalanceID(raw xdr.BalanceId) int64
}

type balanceProvider interface {
	// GetBalanceByID returns history balance struct for specific balance id
	GetBalanceByID(balanceID xdr.BalanceId) history2.Balance
}

type requestProvider interface {
	// GetInvoiceRequestsByContractID returns invoice request which attached to contract with specific id
	GetInvoiceRequestsByContractID(contractID int64) []xdr.ReviewableRequestEntry
}

// OperationHandler used to get info and changes from success operation
type OperationHandler interface {
	// OperationDetails returns db suitable operation details,
	// returns error if operation has not existing action (union switch)
	OperationDetails(op RawOperation, opRes xdr.OperationResultTr) (history2.OperationDetails, error)
	// ParticipantsEffects returns slice of participant effects of each participants
	// that was affected by operation, can include effects (changes) on participants balances
	ParticipantsEffects(opBody xdr.OperationBody, opRes xdr.OperationResultTr,
		source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
	) ([]history2.ParticipantEffect, error)
}

// RawOperation - inner struct to pass source with operation body
// as one parameter in OperationDetails method
type RawOperation struct {
	Source xdr.AccountId
	Body   xdr.OperationBody
}

func customDetailsUnmarshal(rawDetails []byte) map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(rawDetails), &result)
	if err != nil {
		result = make(map[string]interface{})
		result["data"] = string(rawDetails)
		result["error"] = err.Error()
	}

	return result
}

// TODO set option operation handler

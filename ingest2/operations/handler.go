package operations

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	core "gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/generator"
	"gitlab.com/tokend/horizon/ingest2/internal"
	"gitlab.com/tokend/horizon/log"
)

//go:generate mockery -case underscore -name operationsStorage -inpkg -testonly
type operationsStorage interface {
	// Insert - saves to storage operations
	Insert(ops []history2.Operation) error
}

//go:generate mockery -case underscore -name participantEffectsStorage -inpkg -testonly
type participantEffectsStorage interface {
	// Insert - saves to storage operation participant effects
	Insert(participants []history2.ParticipantEffect) error
}

//Handler - handles txs to create operation details and participant effects. Routes operation
// to particular implementation of handler
type Handler struct {
	participantEffectsStorage participantEffectsStorage
	operationsStorage         operationsStorage
	allHandlers               map[xdr.OperationType]handler
	pubKeyProvider            IDProvider
}

// NewOperationsHandler returns new handler which can return
// details and participants effects of certain operation
func NewOperationsHandler(operationsStorage operationsStorage, participantEffectsStorage participantEffectsStorage,
	pubKeyProvider IDProvider, balanceProvider balanceProvider) *Handler {
	manageOfferOpHandlerInst := &manageOfferOpHandler{
		pubKeyProvider: pubKeyProvider,
	}
	return &Handler{
		pubKeyProvider:            pubKeyProvider,
		participantEffectsStorage: participantEffectsStorage,
		operationsStorage:         operationsStorage,
		allHandlers: map[xdr.OperationType]handler{
			xdr.OperationTypeCreateAccount: &createAccountOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeManageAccount: &manageAccountOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeManageExternalSystemAccountIdPoolEntry: &manageExternalSystemPoolOpHandler{},
			xdr.OperationTypeBindExternalSystemAccountId:            &bindExternalSystemAccountOpHandler{},
			xdr.OperationTypeManageBalance: &manageBalanceOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeManageKeyValue: &manageKeyValueOpHandler{},
			xdr.OperationTypeManageLimits: &manageLimitsOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeManageAsset:     &manageAssetOpHandler{},
			xdr.OperationTypeManageAssetPair: &manageAssetPairOpHandler{manageOfferOpHandlerInst},
			xdr.OperationTypeManageOffer:     manageOfferOpHandlerInst,
			xdr.OperationTypeSetFees: &setFeeOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeCreateChangeRoleRequest: &createChangeRoleRequestOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeCreatePreissuanceRequest: &createPreIssuanceRequestOpHandler{},
			xdr.OperationTypeCreateIssuanceRequest: &createIssuanceRequestOpHandler{
				balanceProvider: balanceProvider,
			},
			xdr.OperationTypeCreateSaleRequest: &createSaleRequestOpHandler{},
			xdr.OperationTypeCreateAswapBidRequest: &createAtomicSwapBidRequestOpHandler{
				balanceProvider: balanceProvider,
			},
			xdr.OperationTypeCreateAswapRequest: &createAtomicSwapRequestOpHandler{},
			xdr.OperationTypeCreateWithdrawalRequest: &createWithdrawRequestOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeCreateAmlAlert: &createAMLAlertReqeustOpHandler{
				balanceProvider: balanceProvider,
			},
			xdr.OperationTypeCreateManageLimitsRequest: &createManageLimitsRequestOpHandler{},
			xdr.OperationTypeReviewRequest:             newReviewRequestOpHandler(pubKeyProvider, balanceProvider),
			xdr.OperationTypePaymentV2: &paymentOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeCheckSaleState: &checkSaleStateOpHandler{
				manageOfferOpHandler: manageOfferOpHandlerInst,
			},
			xdr.OperationTypeCancelAswapBid: &cancelAtomicSwapBidOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeSetOptions:           &stubOpHandler{},
			xdr.OperationTypeManageInvoiceRequest: &deprecatedOpHandler{},
			xdr.OperationTypeManageSale: &manageSaleHandler{
				manageOfferOpHandler: manageOfferOpHandlerInst,
			},
			xdr.OperationTypeManageContractRequest: &deprecatedOpHandler{},
			xdr.OperationTypeManageContract:        &deprecatedOpHandler{},
			xdr.OperationTypeCancelSaleRequest:     &stubOpHandler{},
			xdr.OperationTypePayout:                &payoutHandler{},
			xdr.OperationTypeManageAccountRole:     &manageAccountRoleOpHandler{},
			xdr.OperationTypeManageAccountRule:     &manageAccountRuleOpHandler{},
		},
	}
}

// Handle - processes all participants for specific ledger
func (h *Handler) Handle(header *core.LedgerHeader, txs []core.Transaction) error {
	var ledgerOperations []history2.Operation
	var ledgerParticipants []history2.ParticipantEffect
	txIDGen := generator.NewIDI32(header.Sequence)
	opIDGen := generator.NewIDI32(header.Sequence)
	participantEffectIDGen := generator.NewIDI32(header.Sequence)
	for txI := range txs {
		tx := txs[txI]
		txID := txIDGen.Next()
		ops := tx.Envelope.Tx.Operations
		for opI := range ops {
			opDetails, participants, err := h.convertOperation(operation{
				tx:  tx,
				opI: opI,
			}, opIDGen, participantEffectIDGen)
			if err != nil {
				return errors.Wrap(err, "failed to process ledger change", log.F{
					"ledger_seq": header.Sequence,
					"tx_i":       txI,
					"op_i":       opI,
				})
			}

			opDetails.TxID = txID
			opDetails.LedgerCloseTime = internal.TimeFromXdr(xdr.Uint64(header.CloseTime))
			ledgerOperations = append(ledgerOperations, opDetails)
			ledgerParticipants = append(ledgerParticipants, participants...)
		}
	}

	err := h.operationsStorage.Insert(ledgerOperations)
	if err != nil {
		return errors.Wrap(err, "failed to insert operations for ledger", logan.F{
			"ledger_seq": header.Sequence,
			"len(ops)":   len(ledgerOperations),
		})
	}

	err = h.participantEffectsStorage.Insert(ledgerParticipants)
	if err != nil {
		return errors.Wrap(err, "failed to insert operation participants for ledger", logan.F{
			"ledger_seq":    header.Sequence,
			"len(particip)": len(ledgerParticipants),
		})
	}

	return nil
}

// convertOperation transforms xdr operation data to db suitable Operation and Participants Effects
func (h *Handler) convertOperation(op operation, opIDGen *generator.ID,
	participantIDGen *generator.ID) (history2.Operation, []history2.ParticipantEffect, error) {

	opType := op.Operation().Body.Type
	opHandler, ok := h.allHandlers[opType]
	if !ok {
		return history2.Operation{}, nil, errors.From(
			errors.New("no handler for such operation type"), map[string]interface{}{
				"operation type": opType.String(),
			})
	}

	source := op.Source()
	details, err := opHandler.Details(rawOperation{
		Body:   op.Operation().Body,
		Source: source,
	}, op.Result())
	if err != nil {
		return history2.Operation{}, nil,
			errors.Wrap(err, "failed to get operation details", map[string]interface{}{
				"operation type": int32(opType),
			})
	}

	sourceParticipant := history2.ParticipantEffect{
		AccountID: h.pubKeyProvider.MustAccountID(source),
	}
	participantsEffects, err := opHandler.ParticipantsEffects(op.Operation().Body, op.Result(), sourceParticipant,
		op.LedgerChanges())
	if err != nil {
		return history2.Operation{}, nil,
			errors.Wrap(err, "failed to get participants effects", map[string]interface{}{
				"operation type": int32(opType),
			})
	}

	operationID := opIDGen.Next()
	for i := range participantsEffects {
		participantsEffects[i].OperationID = operationID
		participantsEffects[i].ID = participantIDGen.Next()
	}

	return history2.Operation{
		ID:      operationID,
		Details: details,
		Type:    opType,
		Source:  source.Address(),
	}, participantsEffects, nil
}

//Name - returns name of the handler
func (h *Handler) Name() string {
	return "operation_handler"
}

package operations

// todo maybe rename to effects
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
	effectsBaseHandler := effectsProvider{
		IDProvider:      pubKeyProvider,
		balanceProvider: balanceProvider,
	}
	manageOfferOpHandlerInst := &manageOfferOpHandler{
		effectsProvider: effectsBaseHandler,
	}
	return &Handler{
		pubKeyProvider:            pubKeyProvider,
		participantEffectsStorage: participantEffectsStorage,
		operationsStorage:         operationsStorage,
		allHandlers: map[xdr.OperationType]handler{
			xdr.OperationTypeCreateAccount: &createAccountOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeManageExternalSystemAccountIdPoolEntry: &manageExternalSystemPoolOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeBindExternalSystemAccountId: &bindExternalSystemAccountOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeManageBalance: &manageBalanceOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeManageKeyValue: &manageKeyValueOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeManageLimits: &manageLimitsOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeManageAsset: &manageAssetOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeManageAssetPair: &manageAssetPairOpHandler{manageOfferOpHandlerInst},
			xdr.OperationTypeManageOffer:     manageOfferOpHandlerInst,
			xdr.OperationTypeSetFees: &setFeeOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeCreateChangeRoleRequest: &createChangeRoleRequestOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeCreatePreissuanceRequest: &createPreIssuanceRequestOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeCreateIssuanceRequest: &createIssuanceRequestOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeCreateSaleRequest: &createSaleRequestOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeCreateAswapBidRequest: &createAtomicSwapBidRequestOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeCreateAswapRequest: &createAtomicSwapRequestOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeCreateWithdrawalRequest: &createWithdrawRequestOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeCreateAmlAlert: &createAMLAlertReqeustOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeCreateManageLimitsRequest: &createManageLimitsRequestOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeReviewRequest: newReviewRequestOpHandler(effectsBaseHandler),
			xdr.OperationTypePayment: &paymentOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeCheckSaleState: &checkSaleStateOpHandler{
				manageOfferOpHandler: manageOfferOpHandlerInst,
			},
			xdr.OperationTypeCancelAswapBid: &cancelAtomicSwapBidOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeManageInvoiceRequest: &deprecatedOpHandler{},
			xdr.OperationTypeManageSale: &manageSaleHandler{
				manageOfferOpHandler: manageOfferOpHandlerInst,
			},
			xdr.OperationTypeManageContractRequest: &deprecatedOpHandler{},
			xdr.OperationTypeManageContract:        &deprecatedOpHandler{},
			xdr.OperationTypeCancelSaleRequest: &stubOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypePayout: &payoutHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeManageAccountRole: &manageAccountRoleOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeManageAccountRule: &manageAccountRuleOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeManageSigner: &manageSignerOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeManageSignerRole: &manageSignerRoleOpHandler{
				effectsProvider: effectsBaseHandler,
			},
			xdr.OperationTypeManageSignerRule: &manageSignerRuleOpHandler{
				effectsProvider: effectsBaseHandler,
			},
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

	participantsEffects, err := opHandler.ParticipantsEffects(op.Operation().Body, op.Result(), source,
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

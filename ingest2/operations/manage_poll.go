package operations

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
)

type managePollOpHandler struct {
	effectsProvider
}

// Details returns details about manage balance operation
func (h *managePollOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	managePollOp := op.Body.MustManagePollOp()

	details := history2.OperationDetails{
		Type: xdr.OperationTypeManagePoll,
		ManagePoll: &history2.ManagePollDetails{
			Action: managePollOp.Data.Action,
			PollID: int64(managePollOp.PollId),
		},
	}

	switch managePollOp.Data.Action {
	case xdr.ManagePollActionClose:
		closeDetails := &history2.ClosePollData{
			Details:    internal.MarshalCustomDetails(managePollOp.Data.MustClosePollData().Details),
			PollResult: managePollOp.Data.MustClosePollData().Result,
		}
		details.ManagePoll.ClosePoll = closeDetails
	default:
		return history2.OperationDetails{}, errors.From(errors.New("unexpected manage poll action"),
			logan.F{
				"action": managePollOp.Data.Action.ShortString(),
			})
	}

	return details, nil
}

//ParticipantsEffects - returns source of the operation and account for which balance was created if they differ
func (h *managePollOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {

	var participants []history2.ParticipantEffect
	source := h.Participant(sourceAccountID)

	return append(participants, source), nil
}

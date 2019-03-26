package operations

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageVoteOpHandler struct {
	effectsProvider
}

// Details returns details about manage balance operation
func (h *manageVoteOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageVoteOp := op.Body.MustManageVoteOp()

	details := history2.OperationDetails{
		Type: xdr.OperationTypeManageVote,
		ManageVote: &history2.ManageVoteDetails{
			Action: manageVoteOp.Data.Action,
		},
	}

	switch manageVoteOp.Data.Action {
	case xdr.ManageVoteActionCreate:
		details.ManageVote.PollID = int64(manageVoteOp.Data.MustCreateData().PollId)
		voteData := manageVoteOp.Data.MustCreateData().Data
		details.ManageVote.VoteData = &voteData
	case xdr.ManageVoteActionRemove:
		details.ManageVote.PollID = int64(manageVoteOp.Data.MustRemoveData().PollId)
	default:
		return history2.OperationDetails{}, errors.From(errors.New("unexpected manage poll action"),
			logan.F{
				"action": manageVoteOp.Data.Action.ShortString(),
			})
	}

	return details, nil
}

//ParticipantsEffects - returns source of the operation and account for which balance was created if they differ
func (h *manageVoteOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {

	var participants []history2.ParticipantEffect
	source := h.Participant(sourceAccountID)

	return append(participants, source), nil
}

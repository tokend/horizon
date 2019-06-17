package changes

import (
	"encoding/json"
	"time"

	regources "gitlab.com/tokend/regources/generated"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	history "gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
)

type pollStorage interface {
	//Inserts poll into DB
	Insert(poll history.Poll) error
	//Updates poll
	Update(poll history.Poll) error
	// SetState - sets state
	SetState(pollID uint64, state regources.PollState) error
}

type pollHandler struct {
	storage pollStorage
}

func newPollHandler(storage pollStorage) *pollHandler {
	return &pollHandler{
		storage: storage,
	}
}

//Created - handles creation of new poll
func (c *pollHandler) Created(lc ledgerChange) error {
	rawPoll := lc.LedgerChange.MustCreated().Data.MustPoll()

	poll, err := c.convertPoll(rawPoll)
	if err != nil {
		return errors.Wrap(err, "failed to convert poll", logan.F{
			"raw_poll":        rawPoll,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	err = c.storage.Insert(poll)
	if err != nil {
		return errors.Wrap(err, "failed to insert poll into DB", logan.F{
			"poll": poll,
		})
	}

	return nil
}

//Removed - handles state of the poll due to it was removed
func (c *pollHandler) Removed(lc ledgerChange) error {

	pollID := uint64(lc.LedgerChange.MustRemoved().MustPoll().Id)
	managePollOp := lc.Operation.Body.MustManagePollOp()
	state, err := c.getPollState(managePollOp)
	if err != nil {
		return errors.Wrap(err, "failed to get poll state")
	}
	err = c.storage.SetState(pollID, state)
	if err != nil {
		return errors.Wrap(err, "failed to set poll state")
	}

	return nil
}

//Updated - handles update of the poll
func (c *pollHandler) Updated(lc ledgerChange) error {
	rawPoll := lc.LedgerChange.MustUpdated().Data.MustPoll()
	poll, err := c.convertPoll(rawPoll)
	if err != nil {
		return errors.Wrap(err, "failed to convert poll ", logan.F{
			"raw_poll":        rawPoll,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	managePollOp := lc.Operation.Body.MustManagePollOp()
	state, err := c.getPollState(managePollOp)
	if err != nil {
		return errors.Wrap(err, "failed to get poll state")
	}
	poll.State = state

	err = c.updatePollDetails(&poll, managePollOp)
	if err != nil {
		return errors.Wrap(err, "failed to merge update poll details")
	}

	err = c.storage.Update(poll)
	if err != nil {
		return errors.Wrap(err, "failed to update poll", logan.F{
			"poll": poll,
		})
	}
	return nil
}

func (c *pollHandler) updatePollDetails(dst *history.Poll, op xdr.ManagePollOp) (err error) {
	switch op.Data.Action {
	case xdr.ManagePollActionClose:
		var existingDetails map[string]json.RawMessage
		if err = json.Unmarshal([]byte(dst.Details), &existingDetails); err != nil {
			return errors.Wrap(err, "failed to unmarshal existing details")
		}

		if err = json.Unmarshal([]byte(op.Data.ClosePollData.Details), &existingDetails); err != nil {
			return errors.Wrap(err, "failed to unmarshal op details")
		}

		dst.Details, err = json.Marshal(existingDetails)
		if err != nil {
			return errors.Wrap(err, "failed to marshal merged details")
		}
	default:
	}
	return
}

func (c *pollHandler) getPollState(op xdr.ManagePollOp) (regources.PollState, error) {
	var state regources.PollState
	switch op.Data.Action {
	case xdr.ManagePollActionCancel:
		state = regources.PollStateCancelled
	case xdr.ManagePollActionUpdateEndTime:
		state = regources.PollStateOpen
	case xdr.ManagePollActionClose:
		switch op.Data.MustClosePollData().Result {
		case xdr.PollResultFailed:
			state = regources.PollStateFailed
		case xdr.PollResultPassed:
			state = regources.PollStatePassed
		default:
			return state, errors.From(errors.New("Unexpected poll result"), logan.F{
				"poll_result": op.Data.MustClosePollData().Result,
			})
		}
	default:
		return state, errors.From(errors.New("Unexpected manage poll action"), logan.F{
			"action": op.Data.Action,
		})
	}
	return state, nil
}

func (c *pollHandler) convertPoll(raw xdr.PollEntry) (history.Poll, error) {

	return history.Poll{
		ID:               int64(raw.Id),
		OwnerID:          raw.OwnerId.Address(),
		ResultProviderID: raw.ResultProviderId.Address(),
		Data: history.PollData{
			Type: raw.Data.Type,
		},
		NumberOfChoices:          uint32(raw.NumberOfChoices),
		PermissionType:           uint32(raw.PermissionType),
		VoteConfirmationRequired: raw.VoteConfirmationRequired,
		StartTime:                time.Unix(int64(raw.StartTime), 0).UTC(),
		EndTime:                  time.Unix(int64(raw.EndTime), 0).UTC(),
		Details:                  internal.MarshalCustomDetails(raw.Details),
		State:                    regources.PollStateOpen,
	}, nil
}

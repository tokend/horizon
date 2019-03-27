package changes

import (
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	history "gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
	"gitlab.com/tokend/regources/v2"
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

	err = c.storage.Insert(*poll)
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
	err := c.storage.SetState(pollID, regources.PollStateClosed)
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

	err = c.storage.Update(*poll)
	if err != nil {
		return errors.Wrap(err, "failed to update poll", logan.F{
			"poll": poll,
		})
	}
	return nil
}

func (c *pollHandler) convertPoll(raw xdr.PollEntry) (*history.Poll, error) {

	return &history.Poll{
		ID:               int64(raw.Id),
		OwnerID:          raw.OwnerId.Address(),
		ResultProviderID: raw.ResultProviderId.Address(),
		Data: regources.PollData{
			Type: raw.Data.Type,
		},
		NumberOfChoices:          uint64(raw.NumberOfChoices),
		PermissionType:           uint64(raw.PermissionType),
		VoteConfirmationRequired: raw.VoteConfirmationRequired,
		StartTime:                time.Unix(int64(raw.StartTime), 0).UTC(),
		EndTime:                  time.Unix(int64(raw.EndTime), 0).UTC(),
		Details:                  internal.MarshalCustomDetails(raw.Details),
		State:                    regources.PollStateOpen,
	}, nil
}

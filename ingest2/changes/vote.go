package changes

import (
	"gitlab.com/tokend/regources/v2"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	history "gitlab.com/tokend/horizon/db2/history2"
)

type voteStorage interface {
	//Inserts vote into DB
	Insert(vote history.Vote) error
	//Updates vote
	Update(vote history.Vote) error

	Remove(voterID string, pollID uint64) error
}

type voteHandler struct {
	storage voteStorage
}

func newVoteHandler(storage voteStorage) *voteHandler {
	return &voteHandler{
		storage: storage,
	}
}

//Created - handles creation of new vote
func (c *voteHandler) Created(lc ledgerChange) error {
	rawVote := lc.LedgerChange.MustCreated().Data.MustVote()

	vote, err := c.convertVote(rawVote)
	if err != nil {
		return errors.Wrap(err, "failed to convert vote", logan.F{
			"raw_vote":        rawVote,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	err = c.storage.Insert(*vote)
	if err != nil {
		return errors.Wrap(err, "failed to insert vote into DB", logan.F{
			"vote": vote,
		})
	}

	return nil
}

//Removed - handles state of the vote due to it was removed
func (c *voteHandler) Removed(lc ledgerChange) error {

	pollID := uint64(lc.LedgerChange.MustRemoved().MustVote().PollId)
	voterID := lc.LedgerChange.MustRemoved().MustVote().VoterId
	err := c.storage.Remove(voterID.Address(), pollID)
	if err != nil {
		return errors.Wrap(err, "failed to remove vote")
	}

	return nil
}

//Updated - handles update of the vote
func (c *voteHandler) Updated(lc ledgerChange) error {
	rawVote := lc.LedgerChange.MustUpdated().Data.MustVote()
	vote, err := c.convertVote(rawVote)
	if err != nil {
		return errors.Wrap(err, "failed to convert vote ", logan.F{
			"raw_vote":        rawVote,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	err = c.storage.Update(*vote)
	if err != nil {
		return errors.Wrap(err, "failed to update vote", logan.F{
			"vote": vote,
		})
	}
	return nil
}

func (c *voteHandler) convertVote(raw xdr.VoteEntry) (*history.Vote, error) {
	choice := uint64(raw.Data.MustSingle().Choice)
	return &history.Vote{
		VoterID: raw.VoterId.Address(),
		PollID:  int64(raw.PollId),
		VoteData: regources.VoteData{
			PollType:     raw.Data.PollType,
			SingleChoice: &choice,
		},
	}, nil
}

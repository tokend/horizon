package storage

import (
	"github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/history2"
)

// ParticipantEffect - helper struct to store `operation participants`
type ParticipantEffect struct {
	repo *db2.Repo
}

// NewOpParticipants - creates new instance of `ParticipantEffect`
func NewOpParticipants(repo *db2.Repo) *ParticipantEffect {
	return &ParticipantEffect{
		repo: repo,
	}
}

//Insert - stores slice of the participant effects in one batch
func (p *ParticipantEffect) Insert(participants []history2.ParticipantEffect) error {
	// TODO: might have issues due to postgres limit on number of params
	query := squirrel.Insert("participant_effects").Columns("id, account_id, balance_id, asset_code, effect, operation_id")
	for _, participant := range participants {
		query = query.Values(participant.ID, participant.AccountID, participant.BalanceID, participant.AssetCode,
			participant.Effect, participant.OperationID)
	}

	_, err := p.repo.Exec(query)
	if err != nil {
		return errors.Wrap(err, "failed to insert participant effects")
	}

	return nil
}

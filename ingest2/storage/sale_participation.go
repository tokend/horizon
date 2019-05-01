package storage

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/history2"
)

// SaleParticipation is helper struct to operate with `sales`
type SaleParticipation struct {
	repo *db2.Repo
}

// NewSaleParticipation - creates new instance of the `SaleParticipation`
func NewSaleParticipation(repo *db2.Repo) *SaleParticipation {
	return &SaleParticipation{
		repo: repo,
	}
}

func convertParticipationToParams(p history2.SaleParticipation) []interface{} {
	return []interface{}{
		p.ID, p.ParticipantID, p.SaleID, p.BaseAmount, p.QuoteAmount, p.BaseAsset, p.QuoteAsset, p.Price,
	}
}

// Insert - inserts new sale
func (q *SaleParticipation) Insert(participations []history2.SaleParticipation) error {
	columns := []string{"id", "participant_id", "sale_id", "base_amount", "quote_amount", "base_asset", "quote_asset", "price"}

	err := saleParticipationBatchInsert(q.repo, participations,
		"sale_participation", columns,
		convertParticipationToParams)
	if err != nil {
		return errors.Wrap(err, "failed to insert sale participations")
	}
	return nil
}

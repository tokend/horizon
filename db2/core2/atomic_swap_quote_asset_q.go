package core2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// AtomicSwapQuoteAssetQ is a helper to aid in configuring queries
// that loads slices or entry of AtomicSwapBidQuoteAsset structs.
type AtomicSwapQuoteAssetQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

func NewAtomicSwapQuoteAssetQ(repo *db2.Repo) AtomicSwapQuoteAssetQ {
	return AtomicSwapQuoteAssetQ{
		repo: repo,
		selector: sq.Select(
			"qa.ask_id",
			"qa.quote_asset",
			"qa.price",
		).From("atomic_swap_quote_asset qa"),
	}
}

func (q AtomicSwapQuoteAssetQ) FilterByIDs(bidIDs []int64) AtomicSwapQuoteAssetQ {
	q.selector = q.selector.Where(sq.Eq{"qa.ask_id": bidIDs})
	return q
}

func (q AtomicSwapQuoteAssetQ) Select() ([]AtomicSwapQuoteAsset, error) {
	var result []AtomicSwapQuoteAsset
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select atomic swap quote assets")
	}

	return result, nil
}

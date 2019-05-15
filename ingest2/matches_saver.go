package ingest2

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	core "gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/generator"
)

type matchesStorage interface {
	Insert(matches []history2.Match) error
}

// MatchesSaver - converts claimed offer atoms into horizon matches and stores them to db
type MatchesSaver struct {
	storage matchesStorage
}

// NewMatchesSaver - creates new instance of MatchesSaver
func NewMatchesSaver(storage matchesStorage) *MatchesSaver {
	return &MatchesSaver{
		storage: storage,
	}
}

// Handle - converts claimed offer atoms into horizon matches and stores them to db
func (h *MatchesSaver) Handle(header *core.LedgerHeader, txs []core.Transaction) error {
	matchIDGen := generator.NewIDI32(header.Sequence)
	opIDGen := generator.NewIDI32(header.Sequence)

	var ledgerMatches []history2.Match

	for _, tx := range txs {
		ops := tx.Envelope.Tx.Operations
		for opI, op := range ops {
			if op.Body.Type != xdr.OperationTypeManageOffer {
				continue
			}

			opID := opIDGen.Next()
			opResult := tx.Result.Result.Result.MustResults()[opI].MustTr().MustManageOfferResult().MustSuccess()

			var opMatches []history2.Match
			var ok bool
			for _, atom := range opResult.OffersClaimed {
				if opMatches, ok = trySquash(opMatches, atom); !ok {
					opMatches = append(opMatches, history2.NewMatch(
						matchIDGen.Next(),
						opID,
						opResult.BaseAsset,
						opResult.QuoteAsset,
						atom,
					))
				}
			}

			ledgerMatches = append(ledgerMatches, opMatches...)
		}
	}

	if len(ledgerMatches) == 0 {
		return nil
	}

	err := h.storage.Insert(ledgerMatches)
	if err != nil {
		return errors.Wrap(err, "failed to insert matches")
	}

	return nil
}

func trySquash(matches []history2.Match, atom xdr.ClaimOfferAtom) (m []history2.Match, ok bool) {
	for i, match := range matches {
		if match.Price == int64(atom.CurrentPrice) {
			match.BaseAmount += int64(atom.BaseAmount)
			match.QuoteAmount += int64(atom.QuoteAmount)

			matches[i] = match

			return matches, true
		}
	}

	return matches, false
}

// Name - returns name of the handler
func (h *MatchesSaver) Name() string {
	return "matches_saver"
}

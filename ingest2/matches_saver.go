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
			opID := opIDGen.Next()
			if op.Body.Type == xdr.OperationTypeManageOffer {
				manageOfferOp := op.Body.MustManageOfferOp()
				manageOfferOpResult := tx.Result.Result.Result.MustResults()[opI].MustTr().MustManageOfferResult().MustSuccess()

				var opMatches []history2.Match
				var ok bool
				for _, atom := range manageOfferOpResult.OffersClaimed {
					if opMatches, ok = trySquash(opMatches, atom); !ok {
						opMatches = append(opMatches, history2.NewMatch(
							matchIDGen.Next(),
							manageOfferOpResult.BaseAsset,
							manageOfferOpResult.QuoteAsset,
							manageOfferOp.OrderBookId,
							opID,
							atom,
						))
					}
				}

				ledgerMatches = append(ledgerMatches, opMatches...)
			}
		}
	}

	if len(ledgerMatches) > 0 {
		err := h.storage.Insert(ledgerMatches)
		if err != nil {
			return errors.Wrap(err, "failed to insert matches")
		}
	}

	return nil
}

func trySquash(matches []history2.Match, atom xdr.ClaimOfferAtom) (m []history2.Match, ok bool) {
	for index, squashedMatch := range matches {
		if squashedMatch.Price == int64(atom.CurrentPrice) {
			squashedMatch.BaseAmount += int64(atom.BaseAmount)
			squashedMatch.QuoteAmount += int64(atom.QuoteAmount)

			matches[index] = squashedMatch

			return matches, true
		}
	}

	return matches, false
}

// Name - returns name of the handler
func (h *MatchesSaver) Name() string {
	return "matches_saver"
}

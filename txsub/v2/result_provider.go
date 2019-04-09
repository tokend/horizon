package txsub

import (
	"bytes"
	"context"
	"encoding/base64"

	"github.com/pkg/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/db2/history2"
)

// ResultsProvider provides transaction submission results by querying the
// connected horizon and stellar core databases.
type ResultsProvider struct {
	Core    *core2.TransactionQ
	History history2.TransactionsQ
}

// ResultByHash implements txsub.ResultProvider
func (rp *ResultsProvider) ResultByHash(ctx context.Context, hash string) *fullResult {
	// query history database
	hr, err := rp.History.GetByHash(hash)
	if err != nil {
		return &fullResult{
			Err: errors.Wrap(err, "failed to get history transaction by hash"),
		}
	}
	if hr != nil {
		return txResultFromHistory(*hr)
	}

	// query core database
	cr, err := rp.Core.GetByHash(hash)
	if err != nil {
		return &fullResult{
			Err: errors.Wrap(err, "failed to get core transaction by hash"),
		}
	}
	if cr != nil {
		return txResultFromCore(*cr)
	}

	// if no result was found in either db, return nil
	return nil
}

func txResultFromHistory(tx history2.Transaction) *fullResult {
	return &fullResult{
		Result: Result{
			Hash:           tx.Hash,
			LedgerSequence: tx.LedgerSequence,
			EnvelopeXDR:    tx.Envelope,
			ResultXDR:      tx.Result,
			ResultMetaXDR:  tx.Meta,
		},
	}
}

func txResultFromCore(tx core2.Transaction) *fullResult {
	// re-encode result to base64
	var raw bytes.Buffer
	_, err := xdr.Marshal(&raw, tx.Result.Result)

	if err != nil {
		return &fullResult{Err: errors.Wrap(err, "Failed to marshal tx result into xdr")}
	}

	trx := base64.StdEncoding.EncodeToString(raw.Bytes())

	// if result is success, send a normal response
	if tx.Result.Result.Result.Code == xdr.TransactionResultCodeTxSuccess {
		return &fullResult{
			Result: Result{
				Hash:           tx.TransactionHash,
				LedgerSequence: tx.LedgerSequence,
				EnvelopeXDR:    tx.MustEnvelopeXDR(),
				ResultXDR:      trx,
				ResultMetaXDR:  tx.MustResultMetaXDR(),
			},
		}
	}

	// if failed, produce a FailedTransactionError
	return &fullResult{
		Err: NewRejectedTxError(trx),
		Result: Result{
			Hash:           tx.TransactionHash,
			LedgerSequence: tx.LedgerSequence,
			EnvelopeXDR:    tx.MustEnvelopeXDR(),
			ResultXDR:      trx,
			ResultMetaXDR:  tx.MustResultMetaXDR(),
		},
	}

}

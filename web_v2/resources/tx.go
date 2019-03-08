package resources

import (
	"strconv"

	"gitlab.com/tokend/horizon/txsub/v2"

	"gitlab.com/tokend/regources/v2"
)

//NewTxKey - creates new Tx Key for specified ID
func NewTxKey(txID int64) regources.Key {
	return regources.Key{
		ID:   strconv.FormatInt(txID, 10),
		Type: regources.TypeTxs,
	}
}

//NewTxKey - creates new Tx Key for specified ID
func NewTxSubmitKey(hash string) regources.Key {
	return regources.Key{
		ID:   hash,
		Type: regources.TypeTxSubmit,
	}
}

func NewTxSuccess(res *txsub.Result) *regources.TxSubmitResponse {
	key := NewTxSubmitKey(res.Hash)
	return &regources.TxSubmitResponse{
		Key: &key,
		TxSubmitSuccess: &regources.TxSubmitSuccess{
			Data: regources.TxSubmitSuccessData{
				LedgerSequence: res.LedgerSequence,
				Envelope:       res.EnvelopeXDR,
				ResultXDR:      res.ResultXDR,
				Meta:           res.ResultMetaXDR,
			},
		},
	}
}

func NewTxFailure(env txsub.EnvelopeInfo, err txsub.Error) *regources.TxSubmitResponse {
	key := NewTxSubmitKey(env.ContentHash)
	return &regources.TxSubmitResponse{
		Key: &key,
		TxSubmitFailure: &regources.TxSubmitFailure{
			Errors: regources.TxSubmitError{
				Status: err.Status(),
				Detail: err.Details(),
				Title:  err.Type().String(),
				Meta: regources.TxSubmitErrorMeta{
					Envelope:  env.RawBlob,
					ResultXDR: err.ResultXDR(),
				},
			},
		},
	}
}

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
	return &regources.TxSubmitResponse{
		TxSubmitSuccess: &regources.TxSubmitSuccess{
			Key: NewTxSubmitKey(res.Hash),
			Attributes: regources.TxSubmitSuccessAttributes{
				LedgerSequence: res.LedgerSequence,
				Envelope:       res.EnvelopeXDR,
				ResultXDR:      res.ResultXDR,
				Meta:           res.ResultMetaXDR,
			},
		},
	}
}

func NewTxFailure(env txsub.EnvelopeInfo, err txsub.Error) *regources.TxSubmitResponse {
	return &regources.TxSubmitResponse{
		TxSubmitFailure: &regources.TxSubmitFailure{
			Status: err.Status(),
			Detail: err.Details(),
			Title:  err.Type().String(),
			Meta: regources.TxSubmitErrorMeta{
				Envelope:  env.RawBlob,
				ResultXDR: err.ResultXDR(),
			},
		},
	}
}

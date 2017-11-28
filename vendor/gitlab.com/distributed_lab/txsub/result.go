package txsub

// Result - struct providing with all required info to consider result of submission
type Result struct {
	Err            error
	Hash           string
	LedgerSequence int32
	EnvelopeXDR    string
	ResultXDR      string
	ResultMetaXDR  string
}

// HasInternalError - returns true if error is internal
func (r *Result) HasInternalError() bool {
	if r.Err == nil {
		return false
	}

	_, isTxError := r.Err.(Error)
	return !isTxError
}

package txsub

// Result - struct providing with all required info to consider result of submission
type Result struct {
	Hash           string
	LedgerSequence int32
	EnvelopeXDR    string
	ResultXDR      string
	ResultMetaXDR  string
}

type fullResult struct {
	Result
	Err error
}

func (r fullResult) unwrap() (*Result, error) {
	return &r.Result, r.Err
}

// HasInternalError - returns true if error is internal
func (r *fullResult) hasInternalError() bool {
	if r.Err == nil {
		return false
	}

	_, isTxError := r.Err.(Error)
	return !isTxError
}

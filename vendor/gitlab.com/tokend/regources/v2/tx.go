package regources

type TxSubmitResponse struct {
	*Key
	*TxSubmitSuccess
	*TxSubmitFailure
}

type TxSubmitSuccess struct {
	Data TxSubmitSuccessData `json:"data"`
}

type TxSubmitSuccessData struct {
	LedgerSequence int32  `json:"ledger"`
	Envelope       string `json:"envelope_xdr"`
	ResultXDR      string `json:"result_xdr"`
	Meta           string `json:"result_meta_xdr"`
}

type TxSubmitFailure struct {
	Errors TxSubmitError `json:"errors"`
}

type TxSubmitError struct {
	Status int               `json:"status"`
	Title  string            `json:"title"`
	Detail string            `json:"detail"`
	Meta   TxSubmitErrorMeta `json:"meta"`
}

type TxSubmitErrorMeta struct {
	Envelope  string `json:"envelope_xdr"`
	ResultXDR string `json:"result_xdr,omitempty"`
	MetaXDR   string `json:"result_meta_xdr,omitempty"`
}

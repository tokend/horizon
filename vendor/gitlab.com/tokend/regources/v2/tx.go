package regources

type TxSubmitResponse struct {
	*TxSubmitSuccess `json:"data"`
	*TxSubmitFailure `json:"errors"`
}

type TxSubmitSuccess struct {
	Key
	Attributes TxSubmitSuccessAttributes `json:"attributes"`
}

type TxSubmitSuccessAttributes struct {
	LedgerSequence int32  `json:"ledger"`
	Envelope       string `json:"envelope_xdr"`
	ResultXDR      string `json:"result_xdr"`
	Meta           string `json:"result_meta_xdr"`
}

type TxSubmitFailure struct {
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

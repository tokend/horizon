package txsub

// EnvelopeInfo provides information needed for tx envelope submission
type EnvelopeInfo struct {
	// Hash of the transaction
	ContentHash string
	// Source account of the tx
	SourceAddress string
	// base64 encoded envelope
	RawBlob string
}

package history

type OfferState int

const (
	OfferStatePending          OfferState = 1 << iota
	OfferStatePartiallyMatched
	OfferStateFullyMatched
	OfferStateCancelled
)

var offerStateStr = map[OfferState]string{
	OfferStatePending:          "pending",
	OfferStatePartiallyMatched: "partially matched",
	OfferStateFullyMatched:     "fully matched",
	OfferStateCancelled:        "cancelled",
}

func (s OfferState) String() string {
	return offerStateStr[s]
}

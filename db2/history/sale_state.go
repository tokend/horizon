package history

type SaleState int

const (
	SaleStateOpen SaleState = 1 << iota
	SaleStateClosed
	SaleStateCanceled
)

var saleStateStr = map[SaleState]string{
	SaleStateOpen:     "open",
	SaleStateClosed:   "closed",
	SaleStateCanceled: "canceled",
}

func (s SaleState) String() string {
	return saleStateStr[s]
}

package history

type SaleState int

const (
	SaleStateOpen SaleState = 1 << iota
	SaleStateClosed
)

var saleStateStr = map[SaleState]string{
	SaleStateOpen:   "open",
	SaleStateClosed: "closed",
}

func (s SaleState) String() string {
	return saleStateStr[s]
}

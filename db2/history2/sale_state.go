package history2

//SaleState - represents state of the sale
type SaleState int

const (
	//SaleStateOpen - sale is open, but might have not started yet
	SaleStateOpen SaleState = 1 << iota
	//SaleStateClosed - sale has been successfully funded
	SaleStateClosed
	//SaleStateCanceled - sale has been canceled by owner, admin or due to not reaching softcap before end date
	SaleStateCanceled
)

var saleStateStr = map[SaleState]string{
	SaleStateOpen:     "open",
	SaleStateClosed:   "closed",
	SaleStateCanceled: "canceled",
}

//String - converts int enum to string
func (s SaleState) String() string {
	return saleStateStr[s]
}

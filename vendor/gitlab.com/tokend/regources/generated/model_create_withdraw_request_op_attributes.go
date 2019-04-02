package regources

type CreateWithdrawRequestOpAttributes struct {
	Amount         Amount  `json:"amount"`
	CreatorDetails Details `json:"creator_details"`
	Fee            Fee     `json:"fee"`
}

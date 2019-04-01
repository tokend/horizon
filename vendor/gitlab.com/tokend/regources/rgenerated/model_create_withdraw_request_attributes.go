package rgenerated

type CreateWithdrawRequestAttributes struct {
	// Amount to be issued
	Amount         Amount  `json:"amount"`
	CreatorDetails Details `json:"creator_details"`
	Fee            Fee     `json:"fee"`
}

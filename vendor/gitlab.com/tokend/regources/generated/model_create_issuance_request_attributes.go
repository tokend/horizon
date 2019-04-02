package regources

type CreateIssuanceRequestAttributes struct {
	// Amount to be issued
	Amount         Amount  `json:"amount"`
	CreatorDetails Details `json:"creator_details"`
}

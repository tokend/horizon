package regources

type CreatePreIssuanceRequestOpAttributes struct {
	Amount         Amount   `json:"amount"`
	CreatorDetails *Details `json:"creator_details,omitempty"`
}

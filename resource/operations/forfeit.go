package operations

type Forfeit struct {
	Base
	Amount      string `json:"amount"`
	ForfeitType int    `json:"forfeit_type"`
}

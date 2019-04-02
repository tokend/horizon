package regources

type ManageBalanceOpRelationships struct {
	Asset              *Relation `json:"asset,omitempty"`
	DestinationAccount *Relation `json:"destination_account,omitempty"`
}

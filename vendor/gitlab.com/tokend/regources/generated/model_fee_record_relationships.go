package regources

type FeeRecordRelationships struct {
	Account     *Relation `json:"account,omitempty"`
	AccountRole *Relation `json:"account_role,omitempty"`
	Asset       *Relation `json:"asset,omitempty"`
}

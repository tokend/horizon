package regources

//CreateAMLAlertRequestOp - details of corresponding op
type CreateAMLAlertRequestOp struct {
	Key
	Attributes    CreateAMLAlertRequestOpAttrs     `json:"attributes"`
	Relationships CreateAMLAlertRequestOpRelations `json:"relationships"`
}

//CreateAMLAlertRequestOpAttrs - details of corresponding op
type CreateAMLAlertRequestOpAttrs struct {
	Amount Amount `json:"amount"`
	Reason string `json:"reason"`
}

//CreateAMLAlertRequestOpRelations - relationships ot the operation
type CreateAMLAlertRequestOpRelations struct {
	Balance *Relation `json:"balance"`
}

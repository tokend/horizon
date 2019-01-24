package regources

import "time"

// Operation - represent operation
type Operation struct {
	Key
	Relationships OperationRelation `json:"relationships"`
	Attributes    OperationAttr     `json:"attributes"`
}

//OperationRelation - represents operation relationships
type OperationRelation struct {
	Tx     *Relation `json:"tx"`
	Source *Relation `json:"source"`
}

//OperationAttr - represents attributes of operation
type OperationAttr struct {
	Details   OperationDetails `json:"details"`
	AppliedAt time.Time        `json:"applied_at"`
}

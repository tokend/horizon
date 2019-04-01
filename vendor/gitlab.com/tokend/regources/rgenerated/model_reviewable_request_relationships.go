package rgenerated

type ReviewableRequestRelationships struct {
	RequestDetails *Relation `json:"request_details,omitempty"`
	Requestor      *Relation `json:"requestor,omitempty"`
	Reviewer       *Relation `json:"reviewer,omitempty"`
}
